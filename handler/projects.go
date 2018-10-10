package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/gosimple/slug"
	"github.com/labstack/echo"
)

type projectsPageData struct {
	CSRFToken string
	Balances  []store.ProjectBalance
}

func (h Handler) Projects(c echo.Context) error {
	balances, err := h.store.ProjectsBalances()
	if err != nil {
		return errors.New("failed to get project balances from DB: " +
			err.Error())
	}

	return c.Render(http.StatusOK, "projects", projectsPageData{
		CSRFToken: c.Get("csrf-token").(string),
		Balances:  balances,
	})
}

type projectEditPageData struct {
	CSRFToken string
	Project   store.Project
}

func (h Handler) ProjectEdit(c echo.Context) error {
	id := c.Param("project-id")

	project, found, err := h.store.ProjectGet(id)
	if err != nil {
		return errors.New("failed to get project from DB: " + err.Error())
	}
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "project not found")
	}

	return c.Render(http.StatusOK, "project/edit", projectEditPageData{
		CSRFToken: c.Get("csrf-token").(string),
		Project:   project,
	})
}

type projectUsersPageData struct {
	Project  store.Project
	Balances []store.UserBalance
}

func (h Handler) ProjectUsers(c echo.Context) error {
	projectID := c.Param("project-id")

	project, found, err := h.store.ProjectGet(projectID)
	if err != nil {
		return errors.New("failed to get project from DB: " + err.Error())
	}
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "project not found")
	}

	balances, err := h.store.ProjectUsersBalances(projectID)
	if err != nil {
		return errors.New("failed to get project users balances from DB: " +
			err.Error())
	}

	return c.Render(http.StatusOK, "project/users", projectUsersPageData{
		Project:  project,
		Balances: balances,
	})
}

func (h Handler) NewProject(c echo.Context) error {
	name := c.FormValue("name")
	name = strings.TrimSpace(name)
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "blank name")
	}

	id := slug.Make(name)

	err := h.store.ProjectAdd(id, name)
	if err != nil {
		return errors.New("failed to add project to DB: " + err.Error())
	}

	return c.Redirect(http.StatusFound, "/projects")
}

var (
	newNameRe = regexp.MustCompile(`\S`)
)

func (h Handler) EditProject(c echo.Context) error {
	id := c.Param("project-id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "blank project ID")
	}

	action := c.FormValue("action")

	switch action {
	case "edit":
		newName := c.FormValue("name")
		newName = strings.TrimSpace(newName)

		if !newNameRe.MatchString(newName) {
			return echo.NewHTTPError(http.StatusBadRequest,
				"invalid project name")
		}

		err := h.store.ProjectSet(id, newName)
		if err != nil {
			return errors.New("failed to set in DB: " + err.Error())
		}

		return c.Redirect(http.StatusFound, "/projects/"+id+"/edit")

	case "remove":
		err := h.store.ProjectRemove(id)
		if err != nil {
			return errors.New("failed to remove from DB: " + err.Error())
		}

		return c.Redirect(http.StatusFound, "/projects")
	}

	return echo.NewHTTPError(http.StatusBadRequest, "unknown action")
}
