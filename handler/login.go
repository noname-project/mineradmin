package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/boomstarternetwork/bestore"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type loginPageData struct {
	CSRFToken string
	Path      string
}

func (h Handler) Login(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		path := c.QueryParam("path")
		if len(path) > 0 && path[0] != '/' {
			path = ""
		}
		return c.Render(http.StatusOK, "login", loginPageData{
			CSRFToken: c.Get("csrf-token").(string),
			Path:      path,
		})
	}

	login := c.FormValue("login")
	password := c.FormValue("password")
	path := c.FormValue("path")

	err := h.store.CheckAdminPassword(login, password)
	if err != nil {
		if bestore.InvalidLoginOrPassword(err) {
			return echo.NewHTTPError(http.StatusBadRequest,
				"invalid login or password")
		}
		return errors.New("failed to check password in DB: " + err.Error())
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["login"] = login
	claims["exp"] = time.Now().Add(time.Hour * 12).Unix()

	tokenEnc, err := token.SignedString(h.jwtSecret)
	if err != nil {
		return errors.New("failed to sign authorization token: " + err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = tokenEnc
	cookie.Expires = time.Now().Add(time.Hour * 12)

	c.SetCookie(cookie)

	if path == "" {
		path = "/"
	}

	c.Logger().Info("path: ", path)

	return c.Redirect(http.StatusFound, path)
}

func (h Handler) Logout(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "auth"
	cookie.Value = ""
	cookie.Expires = time.Time{}

	c.SetCookie(cookie)

	return c.Redirect(http.StatusFound, "/login")
}
