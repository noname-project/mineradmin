package handler

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func RegisterRenderer(e *echo.Echo) {
	t, err := template.New("projects").Parse(projectsTemplate)
	if err != nil {
		e.Logger.Fatal(err)
	}

	t, err = t.New("project-users").Parse(projectUsersTemplate)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.Renderer = &Template{templates: t}
}
