package handler

import (
	"net/http"

	"github.com/boomstarternetwork/bestore"
	"github.com/labstack/echo"
)

type Handler struct {
	store     bestore.Store
	jwtSecret []byte
}

func NewHandler(s bestore.Store, jwtSecret string) Handler {
	return Handler{
		store:     s,
		jwtSecret: []byte(jwtSecret),
	}
}

func (h Handler) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index", nil)
}
