package service

import (
	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
	"gorm.io/gorm"
)

// AccountService is a service for managing accounts in the database
// Instances of this service should be created using the NewAccountService function
type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{
		db: db,
	}
}

// TODO: Maybe move this to a separate file (own package?)
type IAccountService interface {
	FindAccounts(dest interface{}, meta *pagination.Meta) *gorm.DB
	FindAccountsWithTodos(dest interface{}, conditions ...interface{}) *gorm.DB
	FindAccount(dest interface{}, conditions ...interface{}) *gorm.DB
	FindAccountByID(dest interface{}, id uint) *gorm.DB
	FindAccountByEmail(dest interface{}, email string) *gorm.DB

	CreateAccount(account *model.Account) *gorm.DB
}

// TODO: Maybe cleanup the base model call
func (as *AccountService) FindAccounts(dest interface{}, meta *pagination.Meta) *gorm.DB {
	return model.FindWithMeta(as.db, dest, &model.Account{}, meta, nil)
}

func (as *AccountService) FindAccountsWithTodos(dest interface{}, conditions ...interface{}) *gorm.DB {
	return as.db.Model(&model.Account{}).Preload("Todos").Find(dest, conditions...)
}

func (as *AccountService) FindAccount(dest interface{}, conditions ...interface{}) *gorm.DB {
	return as.db.Model(&model.Account{}).Take(dest, conditions...)
}

// TODO: This was a generic function, but I'm not sure how to make it generic again
func (as *AccountService) FindAccountByID(dest interface{}, id uint) *gorm.DB {
	return as.FindAccount(dest, "id = ?", id)
}

// func (as *AccountService) FindAccountByID(dest interface{}, id string) *gorm.DB {
// 	return FindAccount(dest, "id = ?", id)
// }

func (as *AccountService) FindAccountByEmail(dest interface{}, email string) *gorm.DB {
	return as.FindAccount(dest, "email = ?", email)
}

func (as *AccountService) CreateAccount(account *model.Account) *gorm.DB {
	return as.db.Model(&model.Account{}).Create(account)
}
