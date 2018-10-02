package main

import (
	"database/sql"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"bitbucket.org/boomstarternetwork/mineradmin/handler"
	"bitbucket.org/boomstarternetwork/mineradmin/store"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

// Environment variable like
// `"postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full"`
// or
// `user=pqgotest password=password dbname=pqgotest sslmode=verify-full`
const connStrEnv = "MINERADMIN_POSTGRES_CONNECTION_STRING"

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RemoveTrailingSlashWithConfig(middleware.TrailingSlashConfig{
		RedirectCode: http.StatusMovedPermanently,
	}))

	connStr, exists := os.LookupEnv(connStrEnv)
	if !exists {
		e.Logger.Fatalf("%s environment variable haven't set", connStrEnv)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ps := store.NewDBProjectsStore(db)
	bs := store.NewDBBalancesStore(db)

	h := handler.NewHandler(ps, bs)

	t, err := template.New("projects").Parse(handler.ProjectsTemplate)
	if err != nil {
		e.Logger.Fatal(err)
	}

	t, err = t.New("project-users").Parse(handler.ProjectUsersTemplate)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Renderer = &Template{templates: t}

	e.GET("/", h.Index)
	e.GET("/projects", h.Projects)
	e.GET("/project/:project-id/users", h.ProjectUsers)

	e.Logger.Fatal(e.Start(":80"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
