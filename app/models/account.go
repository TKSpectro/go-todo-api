package models

import (
	"crypto/rand"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	BaseModel
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	Firstname   string `gorm:"" json:"firstname"`
	Lastname    string `gorm:"" json:"lastname"`
	SecretToken string `gorm:"type:varchar(8)" json:"secretToken"`

	Todos []Todo
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
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
	ll := len(ALLOWED_SECRET_TOKEN_CHARS)
	// 8 comes from db max length of secretToken
	result := make([]byte, SECRET_TOKEN_MAX_LENGTH)

	rand.Read(result)
	for i := 0; i < SECRET_TOKEN_MAX_LENGTH; i++ {
		result[i] = ALLOWED_SECRET_TOKEN_CHARS[int(result[i])%ll]
	}

	return string(result)
}
