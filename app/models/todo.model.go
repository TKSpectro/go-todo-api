package models

import (
	"github.com/TKSpectro/go-todo-api/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/utils"

	"gorm.io/gorm"
)

type Todo struct {
	BaseModel
	Title       string `gorm:"not null" json:"title" x-search:"true"`
	Description string `gorm:"" json:"description" x-search:"true"`

	AccountID int `gorm:"not null" json:"fkAccountId"`
	// Account   Account
}

func FindTodosByAccount(dest interface{}, meta *pagination.Meta, accountID uint) *gorm.DB {
	return FindWithMeta(dest, &Todo{}, meta, database.DB.Where("account_id = ?", accountID))
}

func FindTodoByID(dest interface{}, id string, accountID uint) *gorm.DB {
	return database.DB.Model(&Todo{}).Where("id = ? AND account_id = ?", id, accountID).Take(dest)
}

func CreateRandomTodo(accountID uint) *gorm.DB {
	return database.DB.Create(&Todo{
		Title:       utils.RandomString(10, ""),
		Description: utils.RandomString(100, ""),
		AccountID:   int(accountID),
	})
}
