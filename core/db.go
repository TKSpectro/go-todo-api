package core

import (
	"tkspectro/vefeast/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

func SetupDatabase() {
	var err error

	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("DB RIP")
	}

	Migrate(DB)
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&model.Account{})
	db.AutoMigrate(&model.Todo{})
}
