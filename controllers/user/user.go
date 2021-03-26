package user

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserController struct {
	db *gorm.DB
}

func New(db *gorm.DB) UserController {
	return UserController{db: db}
}

//read
func (o *UserController) Dashboard(c echo.Context) error {
	return nil
}

func (o *UserController) Login(c echo.Context) error {
	return nil
}

func (o *UserController) Logout(c echo.Context) error {
	return nil
}

//create
func (o *UserController) Register(c echo.Context) error {
	return nil
}

//update
func (o *UserController) ModifyAccount(c echo.Context) error {
	return nil
}

//delete
func (o *UserController) DeleteAccount(c echo.Context) error {
	return nil
}

//list -- might be better placed in /admin
func (o *UserController) ListUsers(c echo.Context) error {
	return nil
}
