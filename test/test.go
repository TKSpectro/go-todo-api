package test

import (
	"fmt"

	"github.com/TKSpectro/go-todo-api/config"
	"github.com/onsi/ginkgo/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB_SERVER_DSN = fmt.Sprintf("%v:%v@tcp(127.0.0.1:%v)/%s", config.DB_USER, config.DB_ROOT_PASSWORD, config.DB_PORT, "")
	DB_DSN        = fmt.Sprintf("%v:%v@tcp(127.0.0.1:%v)/%s", config.DB_USER, config.DB_ROOT_PASSWORD, config.DB_PORT, config.DB_NAME)
)

func New() {
	ginkgo.GinkgoHelper()

	if !config.IS_TEST {
		panic("[New]::IS_TEST is not true")
	}
	if config.DB_NAME != "go_api_test" {
		panic("[New]::Database name is not go_api_test")
	}

	db := db(DB_SERVER_DSN)
	// sqlDB, _ := db.DB()
	// defer sqlDB.Close()

	db.Exec("DROP DATABASE IF EXISTS " + config.DB_NAME)
	db.Exec("CREATE DATABASE IF NOT EXISTS " + config.DB_NAME)

	// TODO: Add atlas migration run here and remove auto migration on app start
}

func ClearTables(tables []string) {
	ginkgo.GinkgoHelper()

	db := db("")
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	clearTables(tables, db)
}

func ClearAllTables() {
	ginkgo.GinkgoHelper()

	db := db("")
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var tables []string
	if err := db.Table("information_schema.tables").Select("table_name").Where("table_schema = ?", config.DB_NAME).Find(&tables).Error; err != nil {
		fmt.Println("[ClearAllTables]::ERROR GETTING TABLES")
		panic(err)
	}

	clearTables(tables, db)
}

func db(dsn string) *gorm.DB {
	ginkgo.GinkgoHelper()

	if dsn == "" {
		dsn = DB_DSN
	}

	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		fmt.Println("[db]::CONNECTION_ERROR")
		panic(err)
	}

	return db
}

func clearTables(tables []string, db *gorm.DB) {
	ginkgo.GinkgoHelper()

	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
