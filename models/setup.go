package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)
var DB *gorm.DB
func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("page-monitor.db"), &gorm.Config{})
	if err != nil {
	  panic("Failed to connect to database!")
	}
	println(database == nil)
	database.AutoMigrate(&Monitor{})
	DB = database
  }

func ConnectTestDatabase() {
	database, err := gorm.Open(sqlite.Open("../page-monitor-test.db"), &gorm.Config{})
	if err != nil {
	  panic("Failed to connect to database!")
	}
	database.AutoMigrate(&Monitor{})
	DB = database
  }