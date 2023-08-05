package test

import (
	"fmt"
	"log"

	"github.com/TKSpectro/go-todo-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Setup() {
	fmt.Println("[TEST]::SETUP")
	config.IS_TEST = true

	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:%v)/%s", config.DB_USER, config.DB_ROOT_PASSWORD, config.DB_PORT, "")

	fmt.Println("[DATABASE]::CONNECTING", dsn)

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		log.Println("[DATABASE]::CONNECTION_ERROR")
		panic(err)
	}

	fmt.Println("[DATABASE]::CONNECTED", config.DB_NAME)

	if config.DB_NAME != "go_api_test" {
		panic("Database name is not go_api_test")
	}

	fmt.Println("[DATABASE]::RECREATING")
	db.Exec("DROP DATABASE IF EXISTS " + config.DB_NAME)
	db.Exec("CREATE DATABASE IF NOT EXISTS " + config.DB_NAME)

	// TODO: Add atlas migration run here and remove auto migration on app start

	sqlDB, _ := db.DB()
	sqlDB.Close()
}
