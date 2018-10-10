package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boomstarternetwork/mineradmin/mockstore"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Index(t *testing.T) {
	e := echo.New()
	e.Renderer = testRenderer{}

	s := mockstore.New()

	api := NewHandler(s, "secret")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	if assert.NoError(t, api.Index(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}

	s.AssertExpectations(t)
}
