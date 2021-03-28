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

//read
func (o *Controller) Dashboard(c echo.Context) error {
	return c.Render(http.StatusOK, "dashboard.html", nil)
}

func (o *Controller) Login(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}

func (o *Controller) Logout(c echo.Context) error {
	return c.Render(http.StatusOK, "logout.html", nil)
}

func (o *Controller) Forgot(c echo.Context) error {
	return c.Render(http.StatusOK, "forgot.html", nil)
}

//create
func (o *Controller) Register(c echo.Context) error {
	return c.Render(http.StatusOK, "register.html", nil)
}

//update
func (o *Controller) EditAccount(c echo.Context) error {
	return c.Render(http.StatusOK, "account.html", nil)
}

//delete
func (o *Controller) DeleteAccount(c echo.Context) error {
	return c.Render(http.StatusOK, "delete.html", nil)
}
