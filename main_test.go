package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/boomstarternetwork/bestore"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	jwtSecret = "secret"
	runMode   = "testing"
	logLevel  = "off"
)

func initTestWebServer() (*bestore.MockStore, *echo.Echo, error) {
	s := bestore.NewMockStore()

	e, err := initWebServer(s, jwtSecret, runMode, logLevel)
	if err != nil {
		return s, e, err
	}

	e.Renderer = testRenderer{}

	return s, e, nil
}

type testRenderer struct{}

func (testRenderer) Render(_ io.Writer, _ string, _ interface{},
	_ echo.Context) error {
	return nil
}

func makeTestingJWTToken() string {
	claims := jwt.MapClaims{
		"login": "login",
		"exp":   time.Now().Add(12 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenEnc, _ := token.SignedString([]byte(jwtSecret))

	return tokenEnc
}

func Test_Projects(t *testing.T) {
	s, e, err := initTestWebServer()
	if !assert.NoError(t, err) {
		return
	}

	s.On("ProjectsBalances").Return([]bestore.ProjectBalance{
		{
			ProjectID:   123,
			ProjectName: "name-1",
			Coins: []bestore.CoinAmount{
				{Coin: bestore.BTC, Amount: "0.1"},
			},
		},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/projects", nil)
	req.AddCookie(&http.Cookie{Name: "auth", Value: makeTestingJWTToken()})

	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	s.AssertExpectations(t)
}

func Test_ProjectEdit(t *testing.T) {
	s, e, err := initTestWebServer()
	if !assert.NoError(t, err) {
		return
	}

	s.On("GetProject", uint(123)).
		Return(bestore.Project{ID: 123, Name: "name"}, nil)

	req := httptest.NewRequest(http.MethodGet, "/projects/123/edit", nil)
	req.AddCookie(&http.Cookie{Name: "auth", Value: makeTestingJWTToken()})

	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	s.AssertExpectations(t)
}

func Test_ProjectUsers(t *testing.T) {
	s, e, err := initTestWebServer()
	if !assert.NoError(t, err) {
		return
	}

	s.On("GetProject", uint(123)).
		Return(bestore.Project{ID: 123, Name: "name"}, nil)

	s.On("ProjectUsersBalances", uint(123)).
		Return([]bestore.UserBalance{
			{
				Email: "email",
				Coins: []bestore.CoinAmount{
					{Coin: bestore.BTC, Amount: "0.1"},
				},
			},
		}, nil)

	req := httptest.NewRequest(http.MethodGet, "/projects/123/users", nil)
	req.AddCookie(&http.Cookie{Name: "auth", Value: makeTestingJWTToken()})

	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	s.AssertExpectations(t)
}

func Test_NewProject(t *testing.T) {
	s, e, err := initTestWebServer()
	if !assert.NoError(t, err) {
		return
	}

	s.On("AddProject", "Test").
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/projects",
		strings.NewReader("name=Test&csrf-token=token"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "auth", Value: makeTestingJWTToken()})
	req.AddCookie(&http.Cookie{Name: "_csrf", Value: "token"})

	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	body, _ := ioutil.ReadAll(res.Body)
	t.Log(string(body))
	assert.Equal(t, http.StatusFound, res.Code)
	assert.Equal(t, "/projects", res.Header().Get("Location"))

	s.AssertExpectations(t)
}

func Test_EditProject_editAction(t *testing.T) {
	s, e, err := initTestWebServer()
	if !assert.NoError(t, err) {
		return
	}

	s.On("SetProjectName", uint(123), "new-name").
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/projects/123",
		strings.NewReader("action=edit&name=new-name&csrf-token=token"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "auth", Value: makeTestingJWTToken()})
	req.AddCookie(&http.Cookie{Name: "_csrf", Value: "token"})

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	body, _ := ioutil.ReadAll(res.Body)
	t.Log(string(body))
	assert.Equal(t, http.StatusFound, res.Code)
	assert.Equal(t, "/projects/123/edit",
		res.Header().Get("Location"))

	s.AssertExpectations(t)
}

func Test_EditProject_removeAction(t *testing.T) {
	s, e, err := initTestWebServer()
	if !assert.NoError(t, err) {
		return
	}

	s.On("RemoveProject", uint(123)).
		Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/projects/123",
		strings.NewReader("action=remove&csrf-token=token"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "auth", Value: makeTestingJWTToken()})
	req.AddCookie(&http.Cookie{Name: "_csrf", Value: "token"})

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res := httptest.NewRecorder()

	e.ServeHTTP(res, req)

	body, _ := ioutil.ReadAll(res.Body)
	t.Log(string(body))
	assert.Equal(t, http.StatusFound, res.Code)
	assert.Equal(t, "/projects", res.Header().Get("Location"))
}
