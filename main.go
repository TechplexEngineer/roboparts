package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
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
	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = t
	e.Static("/static", "static")
	e.GET("/", controllers.Home)

	//r := mux.NewRouter()
	//static := r.PathPrefix("/static/")
	//static.Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	//static.Methods(http.MethodGet)
	//r.HandleFunc("/", controllers.Home)

	//User
	e.GET("/login", user.Login)
	e.POST("/login", user.Login)
	e.GET("/logout", user.Logout)
	e.GET("/register", user.Register)
	e.POST("/register", user.Register)
	e.GET("/leave", user.DeleteAccount)
	e.GET("/dashboard", user.Dashboard)
	e.GET("/account", user.ModifyAccount)
	e.POST("/account", user.ModifyAccount)

	//Project
	e.GET("/projects", project.List)
	e.GET("/project/new", project.Create)
	e.POST("/project/new", project.Create)
	e.GET("/project/:id", project.Read)
	e.POST("/project/:id", project.Update)
	e.DELETE("/project/:id", project.Delete)

	e.Logger.Fatal(e.Start(":8090"))

}
