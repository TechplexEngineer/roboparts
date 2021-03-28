package project

import (
	"github.com/labstack/echo/v4"
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

func (o *Controller) Create(c echo.Context) error {
	return c.Render(http.StatusOK, "create.html", nil)
}
func (o *Controller) Read(c echo.Context) error {
	c.Param("id")
	project := models.Project{}
	o.db.Find(&project, c.Param("id"))
	data := map[string]interface{}{
		"project": project,
	}
	return c.Render(http.StatusOK, "read.html", data)
}
func (o *Controller) Update(c echo.Context) error {
	return c.Render(http.StatusOK, "update.html", nil)
}
func (o *Controller) Delete(c echo.Context) error {
	return c.Render(http.StatusOK, "delete.html", nil)
}
func (o *Controller) List(c echo.Context) error {
	projects := []models.Project{}
	o.db.Find(&projects) //probably want to add a limit
	data := map[string]interface{}{
		"columns": []string{
			"ID",
			"name",
			"part_prefix",
			"archived",
			"notes",
			"parts",
			"orders",
			"CreatedAt",
			"UpdatedAt",
			"DeletedAt",
		},
		"projects": projects,
	}
	return c.Render(http.StatusOK, "list.html", data)
}
