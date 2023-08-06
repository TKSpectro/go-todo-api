package service

import (
	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/utils"
	"gopkg.in/guregu/null.v4/zero"
	"gorm.io/gorm"
)

// /TodoService is a service for managing accounts in the database
// Instances of this service should be created using the NewTodoService function
type TodoService struct {
	db *gorm.DB
}

func NewTodoService(db *gorm.DB) *TodoService {
	return &TodoService{
		db: db,
	}
}

// TODO: Maybe move this to a separate file (own package?)
type ITodoService interface {
	FindTodosByAccount(dest interface{}, meta *pagination.Meta, accountID uint) *gorm.DB
	FindTodoByID(dest interface{}, id string, accountID uint) *gorm.DB
	CreateTodo(todo *model.Todo) *gorm.DB
	UpdateTodo(todo *model.Todo) *gorm.DB
	DeleteTodoByID(id string) *gorm.DB
	CreateRandomTodo(accountID uint) *gorm.DB
}

// TODO: Maybe cleanup the base model call
func (ts *TodoService) FindTodosByAccount(dest interface{}, meta *pagination.Meta, accountID uint) *gorm.DB {
	return model.FindWithMeta(ts.db, dest, &model.Todo{}, meta, ts.db.Where("account_id = ?", accountID))
}

func (ts *TodoService) FindTodoByID(dest interface{}, id string, accountID uint) *gorm.DB {
	return ts.db.Model(&model.Todo{}).Where("id = ? AND account_id = ?", id, accountID).Take(dest)
}

func (ts *TodoService) CreateTodo(todo *model.Todo) *gorm.DB {
	return ts.db.Create(todo)
}

func (ts *TodoService) UpdateTodo(todo *model.Todo) *gorm.DB {
	return ts.db.Save(todo)
}

func (ts *TodoService) DeleteTodoByID(id string) *gorm.DB {
	return ts.db.Delete(&model.Todo{}, id)
}

func (ts *TodoService) CreateRandomTodo(accountID uint) *gorm.DB {
	return ts.db.Create(&model.Todo{
		Title:       zero.StringFrom(utils.RandomString(100, "")),
		Description: zero.StringFrom(utils.RandomString(100, "")),
		AccountID:   accountID,
	})
}
