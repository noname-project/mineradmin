package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Index(t *testing.T) {
	e := echo.New()

	api := NewHandler(store.NewMockProjectsStore(), store.NewMockBalancesStore())

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	if assert.NoError(t, api.Index(c)) {
		assert.Equal(t, http.StatusFound, res.Code)
		assert.Equal(t, "/projects", res.Header().Get("Location"))
	}
}

func TestHandler_Projects(t *testing.T) {
	e := echo.New()

	RegisterRenderer(e)

	ps := store.NewMockProjectsStore()
	bs := store.NewMockBalancesStore()
	bs.On("ProjectsBalances").Return([]store.ProjectBalance{
		{
			ProjectID:   "id-1",
			ProjectName: "name-1",
			Coins: []store.CoinAmount{
				{Coin: "BTC", Amount: "0.1"},
			},
		},
	}, nil)

	h := NewHandler(ps, bs)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	if assert.NoError(t, h.Projects(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}
}

func TestHandler_ProjectUsers(t *testing.T) {
	e := echo.New()

	RegisterRenderer(e)

	ps := store.NewMockProjectsStore()
	ps.On("Get", "id-1").
		Return(store.Project{ID: "id-1", Name: "name-1"}, nil)

	bs := store.NewMockBalancesStore()
	bs.On("ProjectUsersBalances", "id-1").
		Return([]store.UserBalance{
			{
				Address: "address-1",
				Coins: []store.CoinAmount{
					{Coin: "BTC", Amount: "0.1"},
				},
			},
		}, nil)

	h := NewHandler(ps, bs)

	req := httptest.NewRequest(http.MethodGet, "/project/id-1/users", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	c.SetPath("/project/:project-id/users")
	c.SetParamNames("project-id")
	c.SetParamValues("id-1")

	if assert.NoError(t, h.ProjectUsers(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}
}
