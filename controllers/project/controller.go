package project

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
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
	err = c.Validate(proj) //@todo show the errors to the user in a nice way

	// validationData contains a map of field names to error strings
	validationData := map[string]string{}
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Printf("invalidValidationError - %s\n", err)

		} else {
			var valErrs validator.ValidationErrors
			if errors.As(err, &valErrs) {
				for _, err := range valErrs {
					errString := ""
					if err.Tag() == "required" {
						errString = fmt.Sprintf("%s is required", err.Field())
					} else {
						errString = fmt.Sprintf("'%s' failed the '%s' check", err.Field(), err.Tag())
					}

					validationData[err.Field()] = errString
				}
				helpers.SetErrorFlash(c, "Form data is invalid. Please correct the errors below.")
			} else {
				log.Printf("not a validator.ValidationErrors")
			}
		}

		data := map[string]interface{}{
			"project":    proj,
			"validation": validationData,
		}
		// we had validation errors, send them back to the user so they can fix the field.
		return c.Render(http.StatusOK, "create.html", data)
	}

	indent, err := json.MarshalIndent(proj, "", "    ")
	if err != nil {
		return err
	}
	log.Printf("Data: %s", indent)

	data := map[string]interface{}{
		"project": models.Project{},
	}
	return c.Render(http.StatusOK, "create.html", data)
}

func (o *Controller) Read(c echo.Context) error {
	project := models.Project{}
	o.db.Find(&project, c.Param("id"))
	data := map[string]interface{}{
		"project": project,
	}
	return c.Render(http.StatusOK, "read.html", data)
}
func (o *Controller) Update(c echo.Context) error {
	project := models.Project{}
	o.db.Where("id = ?", c.Param("id")).First(&project)
	data := map[string]interface{}{
		"project": project,
	}
	return c.Render(http.StatusOK, "update.html", data)
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
