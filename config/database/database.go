package database

import (
	"fmt"
	"tkspectro/vefeast/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

func Setup() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("[DATABASE]::CONNECTION_ERROR")
		panic(err)
	}

	DB = db

	Migrate(DB)
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Account{})
	db.AutoMigrate(&model.Todo{})
}
