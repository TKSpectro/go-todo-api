package database

import (
	"fmt"
	"time"

	"github.com/TKSpectro/go-todo-api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", config.DB_USER, config.DB_ROOT_PASSWORD, config.DB_PORT, config.DB_NAME)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
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
