package helpers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func writeCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

func readCookie(c echo.Context) (*http.Cookie, error) {
	cookie, err := c.Cookie("username")
	if err != nil {
		return nil, err
	}
	return cookie, nil
}
