package project

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/techplexengineer/gorm-roboparts/helpers"
	"github.com/techplexengineer/gorm-roboparts/models"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Controller struct {
	db *gorm.DB
}

func New(db *gorm.DB) Controller {
	return Controller{db: db}
}

func (o *Controller) CreateGET(c echo.Context) error {
	data := map[string]interface{}{
		"project": models.Project{},
	}
	return c.Render(http.StatusOK, "create.html", data)
}
func (o *Controller) CreatePOST(c echo.Context) error {
	log.Printf("CreatePOST")
	params, err := c.FormParams()
	if err != nil {
		return err
	}
	log.Printf("params: %#v", params)

	var proj models.Project
	err = c.Bind(&proj)
	if err != nil {
		log.Printf("bind error %s", err)
		return err
	}

	if validationData, valid := helpers.Validate(c, &proj); !valid {

		data := map[string]interface{}{
			"project":    proj,
			"validation": validationData,
		}
		// we had validation errors, send them back to the user so they can fix the field.
		return c.Render(http.StatusOK, "create.html", data)
	}

	// gorm does magic to generate uuid and sets created_at and updated_at
	res := o.db.Create(&proj)

	if res.Error != nil {
		helpers.SetErrorFlash(c, "Unable to create Project")
		log.Print(fmt.Errorf("unable to create project - %s", res.Error))
		return c.Render(http.StatusOK, "create.html", map[string]interface{}{
			"project": proj,
		})
	}
	helpers.SetSuccessFlash(c, "Successfully Created Project")
	return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("project", proj.ID))
}

func (o *Controller) ReadGET(c echo.Context) error {
	project := models.Project{}
	o.db.Where("id = ?", c.Param("id")).First(&project)
	data := map[string]interface{}{
		"project": project,
	}
	return c.Render(http.StatusOK, "read.html", data)
}
func (o *Controller) UpdateGET(c echo.Context) error {
	project := models.Project{}
	o.db.Where("id = ?", c.Param("id")).First(&project)
	data := map[string]interface{}{
		"project": project,
	}
	return c.Render(http.StatusOK, "update.html", data)
}

func (o *Controller) UpdatePOST(c echo.Context) error {
	return c.Render(http.StatusOK, "update.html", nil)
}

func (o *Controller) Delete(c echo.Context) error {
	return c.Render(http.StatusOK, "delete.html", nil)
}

func (o *Controller) ListGET(c echo.Context) error {
	var projects []models.Project
	o.db.Find(&projects) //probably want to add a limit
	data := map[string]interface{}{
		"projects": projects,
	}

	return c.Render(http.StatusOK, "list.html", data)
}
