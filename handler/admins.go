package handler

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/labstack/echo"
)

type adminsPageData struct {
	CSRFToken string
	Admins    []store.Admin
}

func (h Handler) Admins(c echo.Context) error {
	admins, err := h.store.AdminsList()
	if err != nil {
		return errors.New("failed to get admins from DB: " + err.Error())
	}
	return c.Render(http.StatusOK, "admins", adminsPageData{
		CSRFToken: c.Get("csrf-token").(string),
		Admins:    admins,
	})
}

var AdminLoginRe = regexp.MustCompile(`[\w._-]*\w[\w._-]*`)

func (h Handler) NewAdmin(c echo.Context) error {
	login := c.FormValue("login")
	login = strings.TrimSpace(login)
	if login == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "blank login")
	}

	if !AdminLoginRe.MatchString(login) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid login format")
	}

	password, err := h.store.AdminAdd(login)
	if err != nil {
		return errors.New("failed to add admin to DB: " + err.Error())
	}

	return c.Render(http.StatusOK, "admin/password", password)
}

func (h Handler) EditAdmin(c echo.Context) error {
	idStr := c.Param("admin-id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid admin ID")
	}

	action := c.FormValue("action")

	switch action {
	case "reset-password":
		newPassword, err := h.store.AdminResetPassword(id)
		if err != nil {
			return errors.New("failed to reset password in DB: " + err.Error())
		}

		return c.Render(http.StatusOK, "admin/password", newPassword)

	case "remove":
		err := h.store.AdminRemove(id)
		if err != nil {
			return errors.New("failed to remove from DB: " + err.Error())
		}

		return c.Redirect(http.StatusFound, "/admins")
	}

	return echo.NewHTTPError(http.StatusBadRequest, "unknown action")
}
