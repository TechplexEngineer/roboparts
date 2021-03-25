package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/techplexengineer/gorm-roboparts/controllers/project"
	"github.com/techplexengineer/gorm-roboparts/controllers/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"

	"github.com/techplexengineer/gorm-roboparts/controllers"
	"github.com/techplexengineer/gorm-roboparts/models"
)

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

	r := mux.NewRouter()
	static := r.PathPrefix("/static/")
	static.Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	static.Methods(http.MethodGet)
	r.HandleFunc("/", controllers.Home)

	//User
	r.HandleFunc("/login", user.Login)
	r.HandleFunc("/logout", user.Logout)
	r.HandleFunc("/register", user.Register)
	r.HandleFunc("/leave", user.DeleteAccount)
	r.HandleFunc("/dashboard", user.Dashboard)
	r.HandleFunc("/user/edit", user.ModifyAccount)

	//Project
	r.HandleFunc("/projects", project.List)
	r.HandleFunc("/project/new", project.Create).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/project/:id", project.Read).Methods(http.MethodGet)
	r.HandleFunc("/project/:id", project.Update).Methods(http.MethodPost)
	r.HandleFunc("/project/:id", project.Delete).Methods(http.MethodDelete)

	err = http.ListenAndServe(":8090", r)
	if err != nil {
		panic(fmt.Errorf("server error - %w", err))
	}

}
