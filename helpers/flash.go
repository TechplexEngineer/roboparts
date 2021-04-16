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
	gob.Register(&FlashMessage{})
}

func SetFlash(c echo.Context, msg string) error {
	return SetFlashMessage(c, FlashMessage{
		Title: "msg",
		Desc:  "",
		Level: "warning",
	})
}
func SetFlashMessage(c echo.Context, msg FlashMessage) error {

	sess, err := session.Get("flash", c)
	if err != nil {
		log.Print("HERE 10")
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
func GetFlashMessages(c echo.Context) []FlashMessage {
	sess, err := session.Get("flash", c)
	if err != nil {
		log.Print("HERE 1")
		return nil
	}

	res := sess.Flashes()
	messages := make([]FlashMessage, len(res))
	log.Printf("HERE 2 -- %d - %#v", len(res), res)
	for i, msg := range res {
		messages[i] = msg.(FlashMessage)
	}

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		log.Print("HERE 3")
		c.Logger().Error(fmt.Errorf("unable to remove flash messages - %w", err))
	}
	return messages

}
