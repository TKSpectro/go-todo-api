package test

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/TKSpectro/go-todo-api/config"
	"github.com/TKSpectro/go-todo-api/pkg/database"
	"github.com/onsi/ginkgo/v2"
	"gorm.io/gorm"
)

func Setup() *gorm.DB {
	ginkgo.GinkgoHelper()

	if !config.IS_TEST {
		panic("[New]::IS_TEST is not true")
	}
	if config.ROOT_PATH == "" {
		panic("[New]::ROOT_PATH is empty. Please set the environment variable GTA_ROOT_PATH to the root path of the project (See makefile:test)")
	}

	os.MkdirAll(config.TEST_FILE_PATH, fs.ModePerm)

	dbServer := database.ConnectToTestServer()

	dbServer.Exec("DROP DATABASE IF EXISTS " + config.TEST_DB_NAME)
	dbServer.Exec("CREATE DATABASE IF NOT EXISTS " + config.TEST_DB_NAME)

	sqlDB, _ := dbServer.DB()
	defer sqlDB.Close()

	db := database.ConnectToTest()

	database.AutoMigrate(db)

	return db
}

func Teardown(db *gorm.DB) {
	ginkgo.GinkgoHelper()

	db.Exec("DROP DATABASE IF EXISTS " + config.TEST_DB_NAME)
	db.Exec("CREATE DATABASE IF NOT EXISTS " + config.TEST_DB_NAME)

	os.RemoveAll(config.TEST_FILE_PATH)

	sqlDB, _ := db.DB()
	defer sqlDB.Close()
}

func ClearTables(db *gorm.DB, tables []string) {
	ginkgo.GinkgoHelper()

	clearTables(tables, db)
}

func ClearAllTables(db *gorm.DB) {
	ginkgo.GinkgoHelper()

	var tables []string
	if err := db.Table("information_schema.tables").Select("table_name").Where("table_schema = ?", config.TEST_DB_NAME).Find(&tables).Error; err != nil {
		fmt.Println("[ClearAllTables]::ERROR GETTING TABLES")
		panic(err)
	}

	clearTables(tables, db)
}

func clearTables(tables []string, db *gorm.DB) {
	ginkgo.GinkgoHelper()

	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, table := range tables {
		db.Exec("TRUNCATE TABLE " + table)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
