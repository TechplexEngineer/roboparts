package main

import (
	"fmt"
	"github.com/techplexengineer/gorm-roboparts/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:root@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
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
		&models.Project{})
	if err != nil {
		panic(fmt.Errorf("migration error - %w", err))
	}

}
