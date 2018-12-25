package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/boomstarternetwork/bestore"
	"github.com/labstack/echo"
)

type projectsPageData struct {
	CSRFToken string
	Balances  []bestore.ProjectBalance
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
	Project   bestore.Project
}

func (h Handler) ProjectEdit(c echo.Context) error {
	idStr := c.Param("project-id")

	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid project ID")
	}

	id := uint(id64)

	project, err := h.store.GetProject(id)
	if err != nil {
		if bestore.NotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "project not found")
		}
		return errors.New("failed to get project from DB: " + err.Error())
	}

	return c.Render(http.StatusOK, "project/edit", projectEditPageData{
		CSRFToken: c.Get("csrf-token").(string),
		Project:   project,
	})
}

type projectUsersPageData struct {
	Project  bestore.Project
	Balances []bestore.UserBalance
}

func (h Handler) ProjectUsers(c echo.Context) error {
	idStr := c.Param("project-id")

	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid project ID")
	}

	id := uint(id64)

	project, err := h.store.GetProject(id)
	if err != nil {
		if bestore.NotFound(err) {
			return echo.NewHTTPError(http.StatusNotFound, "project not found")
		}
		return errors.New("failed to get project from DB: " + err.Error())
	}

	balances, err := h.store.ProjectUsersBalances(id)
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

	err := h.store.AddProject(name)
	if err != nil {
		return errors.New("failed to add project to DB: " + err.Error())
	}

	return c.Redirect(http.StatusFound, "/projects")
}

var (
	newNameRe = regexp.MustCompile(`\S`)
)

func (h Handler) EditProject(c echo.Context) error {
	idStr := c.Param("project-id")

	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid project ID")
	}

	id := uint(id64)

	action := c.FormValue("action")

	switch action {
	case "edit":
		newName := c.FormValue("name")
		newName = strings.TrimSpace(newName)

		if !newNameRe.MatchString(newName) {
			return echo.NewHTTPError(http.StatusBadRequest,
				"invalid project name")
		}

		err := h.store.SetProjectName(id, newName)
		if err != nil {
			return errors.New("failed to set in DB: " + err.Error())
		}

		return c.Redirect(http.StatusFound,
			fmt.Sprintf("/projects/%d/edit", id))

	case "remove":
		err := h.store.RemoveProject(id)
		if err != nil {
			return errors.New("failed to remove from DB: " + err.Error())
		}

		return c.Redirect(http.StatusFound, "/projects")
	}

	return echo.NewHTTPError(http.StatusBadRequest, "unknown action")
}
