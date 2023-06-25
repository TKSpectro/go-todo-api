package models

import (
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Account struct {
	BaseModel
	Email       string `gorm:"uniqueIndex;not null" json:"email" x-search:"true"`
	Password    string `gorm:"not null" json:"-"` // json:"-" means that this field will not be serialized
	Firstname   string `gorm:"" json:"firstname"`
	Lastname    string `gorm:"" json:"lastname"`
	TokenSecret string `gorm:"type:varchar(8)" json:"-"` // json:"-" means that this field will not be serialized

	Todos []Todo `gorm:"foreignKey:AccountID" json:"todos"`
}

func FindAccounts(dest interface{}, conditions ...interface{}) *gorm.DB {
	return database.DB.Model(&Account{}).Find(dest, conditions...)
}

func FindAccountsWithTodos(dest interface{}, conditions ...interface{}) *gorm.DB {
	return database.DB.Model(&Account{}).Preload("Todos").Find(dest, conditions...)
}

func FindAccount(dest interface{}, conditions ...interface{}) *gorm.DB {
	return database.DB.Model(&Account{}).Take(dest, conditions...)
}

func FindAccountByID(dest interface{}, id string) *gorm.DB {
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
