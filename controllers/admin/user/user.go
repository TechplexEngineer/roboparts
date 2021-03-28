package user

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type Controller struct {
	db *gorm.DB
}

func New(db *gorm.DB) Controller {
	return Controller{db: db}
}

//list
func (o *Controller) ListUsers(c echo.Context) error {
	return c.Render(http.StatusOK, "list.html", nil)
}
