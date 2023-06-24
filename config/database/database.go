package database

import (
	"fmt"
	"time"
	"tkspectro/vefeast/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(sqlite.Open(config.DB), &gorm.Config{
		// Logger:  logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time { return time.Now().Local() },
	})
	if err != nil {
		fmt.Println("[DATABASE]::CONNECTION_ERROR")
		panic(err)
	}

	DB = db
}

func Migrate(tables ...interface{}) error {
	return DB.AutoMigrate(tables...)
}
