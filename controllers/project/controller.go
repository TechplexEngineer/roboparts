package project

import (
	"encoding/json"
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
	if err != nil {
		helpers.SetErrorFlash(c, err.Error())
		data := map[string]interface{}{
			"project": proj,
		}
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
		//"columns": []string{
		//	"ID",
		//	"name",
		//	"part_prefix",
		//	"archived",
		//	"notes",
		//	"parts",
		//	"orders",
		//	"CreatedAt",
		//	"UpdatedAt",
		//	"DeletedAt",
		//},
		"projects": projects,
	}
	log.Printf("ListGET Data: %#v", data)
	return c.Render(http.StatusOK, "list.html", data)
}
