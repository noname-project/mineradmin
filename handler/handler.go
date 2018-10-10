package handler

import (
	"net/http"

	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/labstack/echo"
)

type Handler struct {
	store     store.Store
	jwtSecret []byte
}

func NewHandler(s store.Store, jwtSecret string) Handler {
	return Handler{
		store:     s,
		jwtSecret: []byte(jwtSecret),
	}
}

func (h Handler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}
