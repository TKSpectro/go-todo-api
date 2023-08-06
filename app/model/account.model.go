package model

import (
	"github.com/TKSpectro/go-todo-api/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	BaseModel
	Email       string `gorm:"uniqueIndex;not null" json:"email" x-search:"true" validate:"required,email"`
	Password    string `gorm:"not null" json:"-"`
	Firstname   string `gorm:"" json:"firstname" x-search:"true"`
	Lastname    string `gorm:"" json:"lastname" x-search:"true"`
	TokenSecret string `gorm:"type:varchar(8)" json:"-"`

	Todos []Todo `gorm:"foreignKey:AccountID" json:"todos"`
}

func (account *Account) WriteRemote(remote Account) {
	account.Email = remote.Email
	account.Firstname = remote.Firstname
	account.Lastname = remote.Lastname
}

func FindAccounts(dest interface{}, meta *pagination.Meta) *gorm.DB {
	return FindWithMeta(dest, &Account{}, meta, nil)
}

func FindAccountsWithTodos(dest interface{}, conditions ...interface{}) *gorm.DB {
	return database.DB.Model(&Account{}).Preload("Todos").Find(dest, conditions...)
}

func FindAccount(dest interface{}, conditions ...interface{}) *gorm.DB {
	return database.DB.Model(&Account{}).Take(dest, conditions...)
}

func FindAccountByID[T string | uint](dest interface{}, id T) *gorm.DB {
	return FindAccount(dest, "id = ?", id)
}

func FindAccountByEmail(dest interface{}, email string) *gorm.DB {
	return FindAccount(dest, "email = ?", email)
}

func CreateAccount(account *Account) *gorm.DB {
	return database.DB.Create(account)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

const ALLOWED_SECRET_TOKEN_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const SECRET_TOKEN_MAX_LENGTH = 8

func GenerateSecretToken() string {
	return utils.RandomString(SECRET_TOKEN_MAX_LENGTH, ALLOWED_SECRET_TOKEN_CHARS)
}
