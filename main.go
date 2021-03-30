package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"runtime"
	"strings"

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

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	tmpl, err := helpers.LoadBaseTemplates()
	if err != nil {
		return fmt.Errorf("failed to load base templates - %w", err)
	}
	tmpl.Funcs(template.FuncMap{
		// see https://echo.labstack.com/guide/routing/
		// Echo#Reverse(name string, params ...interface{}
		"pathFor": func(name string, params ...interface{}) string {
			uri := c.Echo().Reverse(name, params...)
			if len(uri) < 2 {
				//@todo it would be nice to scan all templates for URLs to see if any of the route names are missing
				log.Warnf("Unable to generate URL for '%s'", name)
			}
			return uri
		},
	})

	cwd, _ := os.Getwd()
	callingPkg := path.Dir(getFrame(2).File)

	file := strings.TrimPrefix(callingPkg, cwd+string(os.PathSeparator))
	file += string(os.PathSeparator) + name

	tmpl, err = tmpl.ParseFiles(file)
	if err != nil {
		log.Printf("failed to load %s templates - %s", file, err)
		return fmt.Errorf("failed to load %s templates - %w", file, err)
	}

	// Add global methods if data is a map
	//if viewContext, isMap := data.(map[string]interface{}); isMap {
	//	viewContext["reverse"] = c.Echo().Reverse
	//	viewContext["routes"] = c.Echo().Routes()
	//}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Printf("error executing template %s - %s", file, err)
	}

	return err
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
	//e.Use(middleware.Logger()) -- hide this for now
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = t
	e.Static("/static", "static")
	e.GET("/", controllers.Home)

	//User
	uc := user.New(db)
	e.GET("/login", uc.Login).Name = "login"
	e.POST("/login", uc.Login).Name = "login_action"
	e.GET("/logout", uc.Logout).Name = "logout"
	e.GET("/register", uc.Register).Name = "register"
	e.POST("/register", uc.Register).Name = "register_action"
	e.GET("/forgot", uc.Forgot).Name = "forgot"
	e.POST("/forgot", uc.Forgot).Name = "forgot_action"
	e.GET("/leave", uc.DeleteAccount).Name = "delete_account"
	e.GET("/dashboard", uc.Dashboard).Name = "user_dashboard"
	e.GET("/account", uc.EditAccount).Name = "edit_account"
	e.POST("/account", uc.EditAccount)

	//Project
	pc := project.New(db)
	e.GET("/projects", pc.List).Name = "projects"
	e.GET("/project/new", pc.Create).Name = "project_new"
	e.POST("/project/new", pc.Create)
	e.GET("/project/:id", pc.Read).Name = "project"
	e.GET("/project/edit/:id", pc.Update).Name = "project_edit"
	e.POST("/project/edit/:id", pc.Update)
	e.GET("/project/delete/:id", pc.Delete).Name = "project_delete"
	e.DELETE("/project/delete/:id", pc.Delete)

	port := os.Getenv("PORT")
	if len(port) < 1 {
		port = "8090"
	}

	db.Create(&models.Part{
		PartNumber:     "R21_Test",
		Type:           "Part",
		Name:           "Blake's Fav Hammer",
		Notes:          "Some Notes",
		Status:         "Status",
		SourceMaterial: "N/A",
		HaveMaterial:   false,
		Quantity:       "4",
		CutLength:      "",
		Priority:       0,
		DrawingCreated: false,
		Project: models.Project{
			Name:       "Robot 2021",
			PartPrefix: "R21",
			Archived:   false,
			Notes:      "Sample Robot Project",
		},
	})

	e.Logger.Fatal(e.Start(":" + port))

}
