package user

import (
	"github.com/labstack/echo/v4"
	"github.com/techplexengineer/gorm-roboparts/helpers"
	"github.com/techplexengineer/gorm-roboparts/models"
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

func (o *Controller) LoginGET(c echo.Context) error {
	return c.Render(http.StatusOK, "login.html", nil)
}
func (o *Controller) LoginPOST(c echo.Context) error {
	user := models.User{}
	result := o.db.Where("username = ?", c.FormValue("username")).Or("email = ?", c.FormValue("username")).Find(&user)
	if result.RowsAffected == 0 {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"flash": []string{
				"User not found",
			},
		})
	}
	// one username matched
	if result.RowsAffected == 1 {
		if helpers.CheckPasswordHash(c.FormValue("password"), user.PasswordHash) {
			// @todo it would be nice to set a flash message about success
			// but echo doesn't have flash built in
			err := helpers.CreateUserSession(c, user)
			if err != nil {
				c.Logger().Error(err)
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}
			_ = helpers.SetSuccessFlash(c, "Successfully logged in")
			// Per MDN HTTP 303 says change VERBS to GET, body will be lost
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections#attr2
			return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("user_dashboard"))
		}
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"flashMessages": []string{
				"Invalid username or password",
			},
		})
	}
	if result.RowsAffected > 1 {
		c.Logger().Errorf("Found more than one user for username: %s", user.Username)
	}
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"flashMessages": []string{
			"Internal server error. Please try again",
		},
	})
}

func (o *Controller) Logout(c echo.Context) error {
	helpers.Logout(c)
	_ = helpers.SetSuccessFlash(c, "Successfully logged out")
	return c.Redirect(http.StatusTemporaryRedirect, c.Echo().Reverse("home"))
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
