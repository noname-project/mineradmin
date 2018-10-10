package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/boomstarternetwork/mineradmin/mockstore"
	"github.com/boomstarternetwork/mineradmin/store"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Projects(t *testing.T) {
	e := echo.New()

	e.Renderer = testRenderer{}

	s := mockstore.New()
	s.On("ProjectsBalances").Return([]store.ProjectBalance{
		{
			ProjectID:   "id-1",
			ProjectName: "name-1",
			Coins: []store.CoinAmount{
				{Coin: "BTC", Amount: "0.1"},
			},
		},
	}, nil)

	h := NewHandler(s, "secret")

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	c.Set("csrf-token", "token")

	if assert.NoError(t, h.Projects(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}

	s.AssertExpectations(t)
}

func TestHandler_ProjectEdit(t *testing.T) {
	e := echo.New()

	e.Renderer = testRenderer{}

	s := mockstore.New()
	s.On("ProjectGet", "id-1").
		Return(store.Project{ID: "id-1", Name: "name-1"}, true, nil)

	h := NewHandler(s, "secret")

	req := httptest.NewRequest(http.MethodGet, "/projects/id-1/edit", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	c.Set("csrf-token", "token")

	c.SetPath("/projects/:project-id/edit")
	c.SetParamNames("project-id")
	c.SetParamValues("id-1")

	if assert.NoError(t, h.ProjectEdit(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}

	s.AssertExpectations(t)
}

func TestHandler_ProjectUsers(t *testing.T) {
	e := echo.New()

	e.Renderer = testRenderer{}

	s := mockstore.New()
	s.On("ProjectGet", "id-1").
		Return(store.Project{ID: "id-1", Name: "name-1"}, true, nil)
	s.On("ProjectUsersBalances", "id-1").
		Return([]store.UserBalance{
			{
				Email: "email-1",
				Coins: []store.CoinAmount{
					{Coin: "BTC", Amount: "0.1"},
				},
			},
		}, nil)

	h := NewHandler(s, "secret")

	req := httptest.NewRequest(http.MethodGet, "/projects/id-1/users", nil)
	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	c.SetPath("/projects/:project-id/users")
	c.SetParamNames("project-id")
	c.SetParamValues("id-1")

	if assert.NoError(t, h.ProjectUsers(c)) {
		assert.Equal(t, http.StatusOK, res.Code)
	}

	s.AssertExpectations(t)
}

func TestHandler_NewProject(t *testing.T) {
	e := echo.New()

	e.Renderer = testRenderer{}

	s := mockstore.New()
	s.On("ProjectAdd", "test", "Test").
		Return(nil)

	h := NewHandler(s, "secret")

	req := httptest.NewRequest(http.MethodPost, "/projects",
		strings.NewReader("name=Test"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	if assert.NoError(t, h.NewProject(c)) {
		assert.Equal(t, http.StatusFound, res.Code)
	}

	s.AssertExpectations(t)
}

func TestHandler_EditProject_editAction(t *testing.T) {
	e := echo.New()

	e.Renderer = testRenderer{}

	s := mockstore.New()
	s.On("ProjectSet", "id-1", "new-name-1").
		Return(nil)

	h := NewHandler(s, "secret")

	req := httptest.NewRequest(http.MethodPost, "/projects/id-1",
		strings.NewReader("action=edit&name=new-name-1"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	c.SetPath("/projects/:project-id")
	c.SetParamNames("project-id")
	c.SetParamValues("id-1")

	if assert.NoError(t, h.EditProject(c)) {
		assert.Equal(t, http.StatusFound, res.Code)
	}

	s.AssertExpectations(t)
}

func TestHandler_EditProject_removeAction(t *testing.T) {
	e := echo.New()

	e.Renderer = testRenderer{}

	s := mockstore.New()
	s.On("ProjectRemove", "id-1").
		Return(nil)

	h := NewHandler(s, "secret")

	req := httptest.NewRequest(http.MethodPost, "/projects/id-1",
		strings.NewReader("action=remove"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	c := e.NewContext(req, res)

	c.SetPath("/projects/:project-id")
	c.SetParamNames("project-id")
	c.SetParamValues("id-1")

	if assert.NoError(t, h.EditProject(c)) {
		assert.Equal(t, http.StatusFound, res.Code)
	}

	s.AssertExpectations(t)
}
