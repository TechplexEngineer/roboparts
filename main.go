package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/techplexengineer/gorm-roboparts/controllers"
	"github.com/techplexengineer/gorm-roboparts/controllers/project"
	"github.com/techplexengineer/gorm-roboparts/controllers/user"
	"github.com/techplexengineer/gorm-roboparts/helpers"
	"github.com/techplexengineer/gorm-roboparts/models"
)

const (
	ReloadTemplates = true
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	tmpl, err := helpers.LoadBaseTemplates()
	if err != nil {
		return fmt.Errorf("failed to load base templates - %w", err)
	}

	file := "controllers/" + name
	tmpl, err = tmpl.ParseFiles(file)
	if err != nil {
		return fmt.Errorf("failed to load %s templates - %w", file, err)
	}

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
		viewContext["routes"] = c.Echo().Routes()
	}

	return tmpl.Execute(w, data)
}

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "user:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//dbfile := "test.db"
	//db, err := gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database - %w", err))
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.COTSPart{},
		&models.Vendor{},
		&models.Order{},
		&models.OrderItem{},
		&models.Part{},
		&models.Project{},
	)
	if err != nil {
		panic(fmt.Errorf("failed to automigrate database - %w", err))
	}

	t := &Template{}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = t
	e.Static("/static", "static")
	e.GET("/", controllers.Home)

	uc := user.New(db)

	//User
	e.GET("/login", uc.Login).Name = "login"
	e.POST("/login", uc.Login)
	e.GET("/logout", uc.Logout).Name = "logout"
	e.GET("/register", uc.Register).Name = "register"
	e.POST("/register", uc.Register)
	e.GET("/leave", uc.DeleteAccount).Name = "delete_account"
	e.GET("/dashboard", uc.Dashboard).Name = "user_dashboard"
	e.GET("/account", uc.ModifyAccount).Name = "modify_account"
	e.POST("/account", uc.ModifyAccount)

	//Project
	e.GET("/projects", project.List).Name = "projects"
	e.GET("/project/new", project.Create).Name = "project_new"
	e.POST("/project/new", project.Create)
	e.GET("/project/:id", project.Read).Name = "project"
	e.POST("/project/:id", project.Update)
	e.DELETE("/project/:id", project.Delete)

	e.Logger.Fatal(e.Start(":8090"))

}
