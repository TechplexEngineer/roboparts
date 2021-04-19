package models

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

func getDb(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := "user:password@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database - %s", err)
	}
	return db
}

// empty all database tables
func cleanDb(t *testing.T, db *gorm.DB) {
	t.Helper()
	var tables []string
	db.Raw("SELECT TABLE_NAME FROM information_schema.tables WHERE table_schema = 'testdb'").Scan(&tables)
	log.Printf("tables: %v", tables)
	db.Exec("SET FOREIGN_KEY_CHECKS=0;")
	for _, table := range tables {
		db.Exec(fmt.Sprintf("truncate table %s;", table))
	}
	db.Exec("SET FOREIGN_KEY_CHECKS=1;")
}

func freshDb(t *testing.T, db *gorm.DB) {
	t.Helper()
	db.Exec("DROP DATABASE testdb;")
	db.Exec("CREATE DATABASE IF NOT EXISTS testdb COLLATE = utf8mb4_general_ci;")
	db.Exec("USE testdb;")
}

func TestModels(t *testing.T) {
	db := getDb(t)
	//cleanDb(t, db)
	freshDb(t, db)

	err := db.AutoMigrate(
		&COTSPart{},
		&Order{},
		&OrderItem{},
		&Part{},
		&Permission{},
		&Project{},
		&Role{},
		&User{},
		&Vendor{},
	)
	if err != nil {
		t.Fatalf("Unable to migrate Project - %s", err)
	}

	db.Create(&User{
		Username:     "techplex",
		Email:        "techplex.engineer@gmail.com",
		PasswordHash: "",
		Roles:        nil,
		Permissions:  nil,
	})

	//db.Create(&Project{
	//	Name:       "Blake",
	//	PartPrefix: "prefix",
	//	Archived:   false,
	//	Notes:      "notes",
	//	Parts:      nil,
	//	Orders:     nil,
	//})
}

//func TestOrderCreate(t *testing.T) {
//	db := getDb(t)
//	//cleanDb(t, db)
//	freshDb(t, db)
//
//	err := db.AutoMigrate(&Order{})
//	if err != nil {
//		t.Fatalf("Unable to migrate Project - %s", err)
//	}
//
//	db.Create(&Order{
//		Status:    "Open",
//		OrderedAt: nil,
//		PaidForBy: "",
//		TaxCost:   0,
//		Notes:     "",
//	})
//}
