package database

import (
	"fmt"
	"log"
	"time"

	"github.com/TKSpectro/go-todo-api/config"
	"github.com/TKSpectro/go-todo-api/pkg/app/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Connect() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", config.DB_USER, config.DB_ROOT_PASSWORD, config.DB_PORT, config.DB_NAME)

	return connect(dsn)
}

// ConnectToServer connects to the database server without specifying a database
func ConnectToServer() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", config.DB_USER, config.DB_ROOT_PASSWORD, config.DB_PORT, "")

	return connect(dsn)
}

// ConnectToTest connects to the test database
func ConnectToTest() *gorm.DB {
	dsn := fmt.Sprintf("%v:%v@tcp(127.0.0.1:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", config.TEST_DB_USER, config.TEST_DB_ROOT_PASSWORD, config.TEST_DB_PORT, config.TEST_DB_NAME)

	return connect(dsn)
}

func connect(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// Logger:  logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time { return time.Now().Local() },
	})
	if err != nil {
		log.Println("[DATABASE]::CONNECTION_ERROR")
		panic(err)
	}

	return db
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Account{}, &model.Todo{})
}

func Disconnect(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("[DATABASE]::DISCONNECTION_ERROR")
		panic(err)
	}

	sqlDB.Close()
}
