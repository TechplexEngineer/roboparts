package helpers

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/techplexengineer/gorm-roboparts/models"
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
