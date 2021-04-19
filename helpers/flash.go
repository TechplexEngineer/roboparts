package helpers

import (
	"encoding/gob"
	"fmt"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"log"
)

type FlashMessage struct {
	Title string
	Desc  string
	// Level should be a valid color for a bootstrap alert
	// docs: https://getbootstrap.com/docs/5.0/components/alerts/#examples
	// Options: [primary, secondary, success, danger, warning, info, light, dark]
	Level string
}

func init() {

	// Session library stores flash messages using gob. So we must register the type for the gob encoding to work
	gob.Register(&FlashMessage{})
}

// SetInfoFlash is a shorthand for setting a generic info message.
func SetInfoFlash(c echo.Context, msg string) error {
	return SetFlashMessage(c, FlashMessage{
		Title: msg,
		Desc:  "",
		Level: "info",
	})
}

// SetSuccessFlash is a shorthand for setting a generic success message.
func SetSuccessFlash(c echo.Context, msg string) error {
	return SetFlashMessage(c, FlashMessage{
		Title: msg,
		Desc:  "",
		Level: "success",
	})
}

// SetSuccessFlash is a shorthand for setting a generic success message.
func SetWarningFlash(c echo.Context, msg string) error {
	return SetFlashMessage(c, FlashMessage{
		Title: msg,
		Desc:  "",
		Level: "warning",
	})
}

// SetErrorFlash is a shorthand for setting a generic success message.
func SetErrorFlash(c echo.Context, msg string) error {
	return SetFlashMessage(c, FlashMessage{
		Title: msg,
		Desc:  "",
		Level: "danger",
	})
}

// SetFlashMessage allows full control of the flash message.
// See the definition of FlashMessage for allowable values of Level
func SetFlashMessage(c echo.Context, msg FlashMessage) error {

	sess, err := session.Get("flash", c)
	if err != nil {
		return fmt.Errorf("unable to get flash session - %w", err)
	}

	sess.AddFlash(msg)

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Print("HERE 11")
		c.Logger().Error(fmt.Errorf("unable to save flash message for user %v - %w", msg, err))
	}
	return nil
}

// GetFlashMessages retrieves a []FlashMessages
// Once retrieved the flash messages are removed from the list of messages to show.
func GetFlashMessages(c echo.Context) []interface{} {
	sess, err := session.Get("flash", c)
	if err != nil {
		log.Print("unable to get session")
		return nil
	}

	res := sess.Flashes()
	//@todo would be nice to return the correct type
	//messages := make([]FlashMessage, len(res))
	//log.Printf("HERE 2 -- %d - %#v", len(res), res)
	//for i, msg := range res {
	//	messages[i] = msg.(FlashMessage)
	//}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Print("HERE 3")
		c.Logger().Error(fmt.Errorf("unable to remove flash messages - %w", err))
	}
	return res

}
