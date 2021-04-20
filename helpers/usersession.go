package helpers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/techplexengineer/gorm-roboparts/models"
	"log"
)

func CreateUserSession(c echo.Context, user models.User) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return fmt.Errorf("unable to get session - %w", err)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, //86400 is seconds in a day
		HttpOnly: true,      //indicates cookie should NOT be accessible via client side script
	}
	sess.Values["username"] = user.Username
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Error(fmt.Errorf("unable to save session for user %s - %w", user.Username, err))
	}
	return nil
}

func GetCurrentUser(c echo.Context) (string, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return "", fmt.Errorf("unable to get session - %w", err)
	}
	user, ok := sess.Values["username"]
	if !ok {
		//log.Print("unable to access username property of session")
		return "", nil
	}

	username, ok := user.(string)
	if !ok {
		log.Print("could not type-assert username to string")
		return "", nil
	}
	return username, nil
}

func IsLoggedIn(c echo.Context) bool {
	user, err := GetCurrentUser(c)
	if err != nil {
		return false
	}
	if len(user) == 0 {
		return false
	}

	return true
}

func Logout(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return fmt.Errorf("unable to get session - %w", err)
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, //86400 is seconds in a day
		HttpOnly: true,      //indicates cookie should NOT be accessible via client side script
	}
	sess.Values["username"] = nil
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Error(fmt.Errorf("unable to save logout - %w", err))
	}
	return nil
}
