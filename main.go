package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/boomstarternetwork/bestore"
	"github.com/boomstarternetwork/mineradmin/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	cli "gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

var (
	jwtAuthError = errors.New("invalid or expired jwt")
)

func main() {
	app := cli.NewApp()
	app.Name = "mineradmin"
	app.Usage = ""
	app.Description = "Miningpool miners admin web server."
	app.Author = "Vadim Chernov"
	app.Email = "v.chernov@boomstarter.ru"
	app.Version = "0.2"

	webServerFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "config file",
			Value: "",
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "postgres-cs",
			Usage: "postgres connection string",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "bind-addr",
			Usage: "web server bind address",
			Value: ":80",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "jwt-secret",
			Usage: "JWT secret used in salting",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "run-mode",
			Usage: "run mode: production or development",
			Value: "production",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:  "log-level",
			Usage: "log level: debug, info, warn, error, off",
			Value: "info",
		}),
	}

	app.Commands = []cli.Command{
		{
			Name:   "web-server",
			Usage:  "run web server",
			Action: webServer,
			Flags:  webServerFlags,
			Before: altsrc.InitInputSourceWithContext(webServerFlags,
				func(c *cli.Context) (altsrc.InputSourceContext, error) {
					config := c.String("config")
					if config != "" {
						return altsrc.NewYamlSourceFromFlagFunc("config")(c)
					}
					return &altsrc.MapInputSource{}, nil
				}),
		},
		{
			Name:   "add-admin",
			Usage:  "add admin to database",
			Action: addAdmin,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "postgres-cs, p",
					Usage: "postgres connection string",
				},
				cli.StringFlag{
					Name:  "login, l",
					Usage: "admin login",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("Error: %v\n", err)
		os.Exit(1)
	}
}

func webServer(c *cli.Context) error {
	connStr := c.String("postgres-cs")
	bindAddr := c.String("bind-addr")
	jwtSecret := c.String("jwt-secret")
	runMode := c.String("run-mode")
	logLevel := c.String("log-level")

	s, err := bestore.NewDBStore(connStr, runMode)
	if err != nil {
		return cli.NewExitError("failed to create new DB store: "+
			err.Error(), 1)
	}

	e, err := initWebServer(s, jwtSecret, runMode, logLevel)
	if err != nil {
		return cli.NewExitError("failed to init web server: "+
			err.Error(), 2)
	}

	err = e.Start(bindAddr)

	return cli.NewExitError("failed to start echo server: "+
		err.Error(), 3)
}

func addAdmin(c *cli.Context) error {
	connStr := c.String("postgres-cs")
	login := c.String("login")

	if !handler.AdminLoginRe.MatchString(login) {
		return errors.New("invalid login format")
	}

	s, err := bestore.NewDBStore(connStr, "production")
	if err != nil {
		return cli.NewExitError("failed to create new DB store: "+
			err.Error(), 5)
	}

	password, err := s.AddAdmin(login)
	if err != nil {
		return errors.New("failed to add admin to DB: " + err.Error())
	}

	fmt.Println("Password:", password)

	return nil
}

func initWebServer(s bestore.Store, jwtSecret string,
	runMode string, logLevel string) (*echo.Echo, error) {
	e := echo.New()

	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:  "form:csrf-token",
		ContextKey:   "csrf-token",
		CookieSecure: true,
	}))

	var err error

	switch runMode {
	case "production":
		e.Use(middleware.Recover())
		e.Renderer, err = handler.NewProdTemplateRenderer("./templates")
		if err != nil {
			return nil, errors.New("failed to create production template " +
				"renderer: " + err.Error())
		}
	case "development":
		e.Use(middleware.Recover())
		e.Debug = true
		e.Renderer = handler.NewDevTemplateRenderer("./templates")
	case "testing":
	default:
		return nil, errors.New("invalid run mode")
	}

	switch logLevel {
	case "debug", "info", "warn", "error":
		e.Use(middleware.Logger())
	}

	switch logLevel {
	case "debug":
		e.Logger.SetLevel(log.DEBUG)
	case "info":
		e.Logger.SetLevel(log.INFO)
	case "warn":
		e.Logger.SetLevel(log.WARN)
	case "error":
		e.Logger.SetLevel(log.ERROR)
	case "off":
	default:
		return nil, errors.New("invalid log level")
	}

	h := handler.NewHandler(s, jwtSecret)

	e.GET("/login", h.Login)
	e.POST("/login", h.Login)

	withJWT := middleware.JWTWithConfig(middleware.JWTConfig{
		ErrorHandler: func(e error) error {
			return jwtAuthError
		},
		SigningKey:    []byte(jwtSecret),
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "login",
		TokenLookup:   "cookie:auth",
	})

	withAuth := func(next echo.HandlerFunc) echo.HandlerFunc {
		return withJWT(func(c echo.Context) error {
			err := next(c)
			if err == jwtAuthError {
				redirectPath := "/login"
				path := c.Request().URL.Path + "?" + c.Request().URL.RawQuery
				if path != "/?" && path != "?" {
					redirectPath += "?path=" + url.QueryEscape(path)
				}
				c.Redirect(http.StatusFound, redirectPath)
			}
			return err
		})
	}

	e.GET("/", withAuth(h.Index))

	e.GET("/logout", withAuth(h.Logout))

	e.GET("/projects", withAuth(h.Projects))
	e.GET("/projects/:project-id/edit", withAuth(h.ProjectEdit))
	e.GET("/projects/:project-id/users", withAuth(h.ProjectUsers))
	e.POST("/projects", withAuth(h.NewProject))
	e.POST("/projects/:project-id", withAuth(h.EditProject))

	e.GET("/admins", withAuth(h.Admins))
	e.POST("/admins", withAuth(h.NewAdmin))
	e.POST("/admins/:admin-id", withAuth(h.EditAdmin))

	e.GET("/users", withAuth(h.Users))
	e.POST("/users", withAuth(h.NewUser))
	e.GET("/users/:user-id/addresses", withAuth(h.UserAddresses))
	e.POST("/users/:user-id/addresses", withAuth(h.EditUserAddresses))

	return e, nil
}
