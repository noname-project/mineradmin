package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/boomstarternetwork/mineradmin/coin"

	"github.com/boomstarternetwork/mineradmin/store"

	"github.com/labstack/echo"
)

type usersPageData struct {
	CSRFToken string
	Users     []store.User
}

func (h Handler) Users(c echo.Context) error {
	users, err := h.store.UsersList()
	if err != nil {
		return errors.New("failed to get users list from DB: " + err.Error())
	}

	return c.Render(http.StatusOK, "users", usersPageData{
		CSRFToken: c.Get("csrf-token").(string),
		Users:     users,
	})
}

func (h Handler) NewUser(c echo.Context) error {
	email := c.FormValue("email")
	email = strings.TrimSpace(email)
	if email == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "blank email")
	}
	if strings.Index(email, "@") == -1 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid email format")
	}

	name := c.FormValue("name")
	name = strings.TrimSpace(name)
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "blank name")
	}

	userID, err := h.store.UserAdd(email, name)
	if err != nil {
		return errors.New("failed to add user to DB: " + err.Error())
	}

	return c.Redirect(http.StatusFound, fmt.Sprintf(
		"/users/%d/addresses", userID))
}

type userAddressesData struct {
	CSRFToken string
	Coins     []string
	User      store.User
	Addresses map[string][]string
}

func (h Handler) UserAddresses(c echo.Context) error {
	userIDStr := c.Param("user-id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	user, found, err := h.store.UserGet(userID)
	if err != nil {
		return errors.New("failed to get user from DB: " + err.Error())
	}
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	addrs, err := h.store.UserAddresses(userID)
	if err != nil {
		return errors.New("failed to get user addresses from DB: " + err.Error())
	}

	return c.Render(http.StatusOK, "user/addresses", userAddressesData{
		CSRFToken: c.Get("csrf-token").(string),
		Coins:     coin.List(),
		User:      user,
		Addresses: addrs,
	})
}

var addressRe = regexp.MustCompile(`\S`)

func (h Handler) EditUserAddresses(c echo.Context) error {
	userIDStr := c.Param("user-id")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid user ID")
	}

	_, found, err := h.store.UserGet(userID)
	if err != nil {
		return errors.New("failed to get user from DB: " + err.Error())
	}
	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}

	action := c.FormValue("action")

	cn := c.FormValue("coin")
	if !coin.Valid(cn) {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid coin")
	}

	address := c.FormValue("address")
	address = strings.TrimSpace(address)
	if !addressRe.MatchString(address) {
		return echo.NewHTTPError(http.StatusBadRequest,
			"invalid address format")
	}

	addrsPath := fmt.Sprintf("/users/%d/addresses", userID)

	switch action {
	case "add":
		err := h.store.UserAddressAdd(userID, cn, address)
		if err != nil {
			return errors.New("failed to add to DB: " + err.Error())
		}

		return c.Redirect(http.StatusFound, addrsPath)

	case "remove":
		err := h.store.UserAddressRemove(userID, cn, address)
		if err != nil {
			return errors.New("failed to remove from DB: " + err.Error())
		}

		return c.Redirect(http.StatusFound, addrsPath)
	}

	return echo.NewHTTPError(http.StatusBadRequest, "unknown action")
}
