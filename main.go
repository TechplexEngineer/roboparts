package main

import (
	"flag"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/wader/gormstore/v2"
	"gorm.io/gorm"

	"github.com/techplexengineer/gorm-roboparts/controllers"
	"github.com/techplexengineer/gorm-roboparts/controllers/project"
	"github.com/techplexengineer/gorm-roboparts/controllers/user"
	"github.com/techplexengineer/gorm-roboparts/helpers"
	"github.com/techplexengineer/gorm-roboparts/models"
)

func main() {

	exampleCfg := flag.Bool("x", false, "print example config to stdout")

	flag.Parse()

	if *exampleCfg {
		json, err := GetExampleConfigJson()
		if err != nil {
			log.Printf("Unable to get example config - %s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", json)
		os.Exit(0)
	}

	cfg, err := GetConfig()
	if err != nil {
		panic(fmt.Errorf("unable to load config - %w", err))
	}

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	var db *gorm.DB
	switch strings.ToLower(cfg.Database.Type) {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.Database.DSN), &gorm.Config{})
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.Database.DSN), &gorm.Config{})
	}
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

	e := echo.New()
	e.HideBanner = true

	// rate is requests per second
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))

	store := gormstore.New(db, cfg.Session.AuthKey, cfg.Session.EncryptionKey)
	e.Use(session.Middleware(store))
	// db cleanup every hour
	// close quit channel to stop cleanup
	quit := make(chan struct{}) //@todo close periodic store cleanup
	go store.PeriodicCleanup(1*time.Hour, quit)

	// Middleware
	//e.Use(middleware.Logger()) // quite verbose
	e.Use(middleware.Recover())

	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = &helpers.TemplateRenderer{}
	e.Static("/static", "static")
	e.GET("/", controllers.Home)

	//User
	uc := user.New(db)
	e.GET("/login", uc.LoginGET).Name = "login"
	e.POST("/login", uc.LoginPOST).Name = "login_action"
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

	pwHash, err := helpers.HashPassword("password")
	if err != nil {
		panic(fmt.Errorf("unable to encrypt password - %w", err))
	}
	db.Create(&models.User{
		Username:     "techplex",
		Email:        "techplex.engineer@gmail.com",
		PasswordHash: pwHash,
	})

	e.Logger.Fatal(e.Start(":" + port))

}
