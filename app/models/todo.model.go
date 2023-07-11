package models

import (
	"github.com/TKSpectro/go-todo-api/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/utils"

	"gopkg.in/guregu/null.v4/zero"

	"gorm.io/gorm"
)

type Todo struct {
	BaseModel
	Title       zero.String `gorm:"not null" json:"title" x-search:"true" swaggertype:"string" validate:"required,min=1"`
	Description zero.String `gorm:"" json:"description" x-search:"true" swaggertype:"string"`
	Completed   bool        `gorm:"default:false" json:"completed"`

	AccountID uint `gorm:"not null" json:"fkAccountId"`
	// Account   Account
}

func (todo *Todo) WriteRemote(remote interface{}) {
	todo.Title = remote.(Todo).Title
	todo.Description = remote.(Todo).Description
	todo.Completed = remote.(Todo).Completed
}

func FindTodosByAccount(dest interface{}, meta *pagination.Meta, accountID uint) *gorm.DB {
	return FindWithMeta(dest, &Todo{}, meta, database.DB.Where("account_id = ?", accountID))
}

func FindTodoByID(dest interface{}, id string, accountID uint) *gorm.DB {
	return database.DB.Model(&Todo{}).Where("id = ? AND account_id = ?", id, accountID).Take(dest)
}

func CreateTodo(todo *Todo) *gorm.DB {
	return database.DB.Create(todo)
}

func UpdateTodo(todo *Todo) *gorm.DB {
	return database.DB.Save(todo)
}

func DeleteTodoByID(id string) *gorm.DB {
	return database.DB.Delete(&Todo{}, id)
}

func CreateRandomTodo(accountID uint) *gorm.DB {
	return database.DB.Create(&Todo{
		Title:       zero.StringFrom(utils.RandomString(100, "")),
		Description: zero.StringFrom(utils.RandomString(100, "")),
		AccountID:   accountID,
	})
}
