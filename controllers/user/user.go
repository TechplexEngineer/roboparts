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
	username := c.FormValue("username")
	password := c.FormValue("password")
	user := models.User{}
	result := o.db.Where("username = ?", username).Or("email = ?", username).Find(&user)
	if result.RowsAffected == 0 {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"flash": []string{
				"User not found",
			},
		})
	}
	// one username matched
	if result.RowsAffected == 1 {
		if helpers.CheckPasswordHash(password, user.PasswordHash) {
			// @todo it would be nice to set a flash message about success
			// but echo doesn't have flash built in
			c.Set("auth", true)
			return c.Redirect(http.StatusTemporaryRedirect, c.Echo().Reverse("user_dashboard"))
		}
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"flash": []string{
				"Invalid username or password",
			},
		})
	}
	if result.RowsAffected > 1 {
		c.Logger().Errorf("Found more than one user for username: %s", username)
	}
	return c.Render(http.StatusOK, "login.html", map[string]interface{}{
		"flash": []string{
			"Internal server error. Please try again",
		},
	})
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
