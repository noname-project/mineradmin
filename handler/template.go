package handler

import (
	"errors"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

type ProdTemplateRenderer struct {
	templates map[string]*template.Template
}

func NewProdTemplateRenderer(templatesPath string) (ProdTemplateRenderer,
	error) {
	ts := map[string]*template.Template{}

	if strings.Index(templatesPath, "./") == 0 {
		templatesPath = templatesPath[2:]
	}

	if templatesPath[len(templatesPath)-1] != '/' {
		templatesPath += "/"
	}

	baseTmpl, err := template.ParseFiles(filepath.Join(
		templatesPath, "layout/base.tmpl"))
	if err != nil {
		return ProdTemplateRenderer{}, err
	}

	tmplFileExtRe := regexp.MustCompile(`\.tmpl$`)

	err = filepath.Walk(templatesPath, func(path string, info os.FileInfo,
		err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !tmplFileExtRe.MatchString(info.Name()) {
			return nil
		}

		name := strings.TrimPrefix(path, templatesPath)
		name = strings.TrimSuffix(name, ".tmpl")

		tmpl, err := baseTmpl.Clone()
		if err != nil {
			return err
		}

		tmpl, err = tmpl.ParseFiles(path)
		if err != nil {
			return err
		}

		ts[name] = tmpl

		return nil
	})
	if err != nil {
		return ProdTemplateRenderer{}, err
	}

	return ProdTemplateRenderer{
		templates: ts,
	}, nil
}

func (t ProdTemplateRenderer) Render(w io.Writer, name string,
	data interface{}, _ echo.Context) error {
	tmpl, exists := t.templates[name]
	if !exists {
		return errors.New("template not found")
	}
	return tmpl.Execute(w, data)
}

type DevTemplateRenderer struct {
	templatesPath string
}

func NewDevTemplateRenderer(templatesPath string) DevTemplateRenderer {
	return DevTemplateRenderer{
		templatesPath: templatesPath,
	}
}

func (r DevTemplateRenderer) Render(w io.Writer, name string, data interface{},
	_ echo.Context) error {
	tmpl, err := template.ParseFiles(
		filepath.Join(r.templatesPath, "layout/base.tmpl"),
		filepath.Join(r.templatesPath, name+".tmpl"))
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

type testRenderer struct{}

func (testRenderer) Render(_ io.Writer, _ string, _ interface{},
	_ echo.Context) error {
	return nil
}
