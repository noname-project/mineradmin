package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/boomstarternetwork/mineradmin/dbstore"
	"github.com/boomstarternetwork/mineradmin/handler"
	xormcore "github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/urfave/cli"
)

const (
	// connStrEnv is environment variable with postgres connection string like
	// `"postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"`
	// or
	// `user=pqgotest password=password dbname=pqgotest sslmode=verify-full`
	connStrEnv = "MINERADMIN_POSTGRES_CONNECTION_STRING"

	// bindAddrEnv is environment variable with bind address like:
	// :80, 127.0.0.1:8080
	bindAddrEnv = "MINERADMIN_BIND_ADDR"

	// modeEnv is environment variable with web server run mode:
	// production or development
	modeEnv = "MINERADMIN_MODE"

	jwtSecretEnv = "MINERADMIN_JWT_SECRET"
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
	app.Version = "0.1"

	app.Commands = []cli.Command{
		{
			Name:    "webserver",
			Aliases: []string{"ws"},
			Usage:   "run webserver",
			Action:  webServer,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "postgres-cs, p",
					Usage:  "postgres connection string",
					EnvVar: connStrEnv,
				},
				cli.StringFlag{
					Name:   "bind-addr, b",
					Usage:  "web server bind address",
					EnvVar: bindAddrEnv,
					Value:  ":80",
				},
				cli.StringFlag{
					Name:   "mode, m",
					Usage:  "run mode: production or development",
					EnvVar: modeEnv,
					Value:  "production",
				},
				cli.StringFlag{
					Name:   "jwt-secret, j",
					Usage:  "JWT secret used in salting",
					EnvVar: jwtSecretEnv,
				},
			},
		},
		{
			Name:    "addadmin",
			Aliases: []string{"aa"},
			Usage:   "add admin to database",
			Action:  addAdmin,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "postgres-cs, p",
					Usage:  "postgres connection string",
					EnvVar: connStrEnv,
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
	runMode := c.String("mode")
	jwtSecret := c.String("jwt-secret")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:  "form:csrf-token",
		ContextKey:   "csrf-token",
		CookieSecure: true,
	}))

	xdb, err := xorm.NewEngine("postgres", connStr)
	if err != nil {
		return cli.NewExitError("failed to open gorm DB: "+err.Error(), 2)
	}
	defer xdb.Close()

	xdb.SetMapper(xormcore.GonicMapper{})
	xdb.Logger().SetLevel(xormcore.LOG_OFF)

	switch runMode {
	case "production":
		e.Logger.SetLevel(log.ERROR)
		e.Renderer, err = handler.NewProdTemplateRenderer("./templates")
		if err != nil {
			return cli.NewExitError(
				"failed to create production template renderer: "+
					err.Error(), 3)
		}
	case "development":
		e.Logger.SetLevel(log.INFO)
		e.Debug = true
		e.Renderer = handler.NewDevTemplateRenderer("./templates")
	default:
		return cli.NewExitError("invalid run mode", 4)
	}

	s := dbstore.New(xdb)
	h := handler.NewHandler(s, jwtSecret)

	e.GET("/login", h.Login)
	e.POST("/login", h.Login)

	r := e.Group("")

	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
		}
	})

	r.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		ErrorHandler: func(e error) error {
			return jwtAuthError
		},
		SigningKey:    []byte(jwtSecret),
		SigningMethod: middleware.AlgorithmHS256,
		ContextKey:    "login",
		TokenLookup:   "cookie:auth",
	}))

	r.GET("/", h.Index)

	r.GET("/logout", h.Logout)

	r.GET("/projects", h.Projects)
	r.GET("/projects/:project-id/edit", h.ProjectEdit)
	r.GET("/projects/:project-id/users", h.ProjectUsers)
	r.POST("/projects", h.NewProject)
	r.POST("/projects/:project-id", h.EditProject)

	r.GET("/admins", h.Admins)
	r.POST("/admins", h.NewAdmin)
	r.POST("/admins/:admin-id", h.EditAdmin)

	r.GET("/users", h.Users)
	r.POST("/users", h.NewUser)
	r.GET("/users/:user-id/addresses", h.UserAddresses)
	r.POST("/users/:user-id/addresses", h.EditUserAddresses)

	return cli.NewExitError("failed to start echo server: "+
		e.Start(bindAddr).Error(), 5)
}

func addAdmin(c *cli.Context) error {
	connStr := c.String("postgres-cs")
	login := c.String("login")

	if !handler.AdminLoginRe.MatchString(login) {
		return errors.New("invalid login format")
	}

	xdb, err := xorm.NewEngine("postgres", connStr)
	if err != nil {
		return cli.NewExitError("failed to open gorm DB: "+err.Error(), 2)
	}
	defer xdb.Close()

	xdb.Logger().SetLevel(xormcore.LOG_OFF)

	s := dbstore.New(xdb)

	password, err := s.AdminAdd(login)
	if err != nil {
		return errors.New("failed to add admin to DB: " + err.Error())
	}

	fmt.Println("Password:", password)

	return nil
}
