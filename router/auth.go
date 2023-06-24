package router

import (
	"crypto/rand"
	"errors"
	"time"
	"tkspectro/vefeast/config/database"
	"tkspectro/vefeast/core"
	"tkspectro/vefeast/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterRoutesAuth(router fiber.Router) {
	// Authentication routes
	router.Put("/login", func(c *fiber.Ctx) error {
		db := database.DB

		remote := new(model.Account)
		if err := c.BodyParser(remote); err != nil {
			return &core.BAD_REQUEST
		}

		var account model.Account
		if err := db.Where(&model.Account{Email: remote.Email}).Find(&account).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return &core.NOT_FOUND
			}
			return err
		}

		if !CheckPasswordHash(remote.Password, account.Password) {
			return &core.UNAUTHORIZED
		}

		// Map out all the claims to write to the payload
		claims := jwt.MapClaims{
			"accountId": account.ID,
			"exp":       time.Now().Add(time.Hour * 72).Unix(),
		}

		// Generate token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// TODO: Switch to RS256 with signing files https://github.com/gofiber/contrib/tree/main/jwt#rs256-example
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return &core.INTERNAL_SERVER_ERROR
		}

		// TODO: Implement refresh token
		return c.JSON(fiber.Map{"token": t})
	})

	router.Post("/register", func(c *fiber.Ctx) error {
		db := database.DB

		account := new(model.Account)
		if err := c.BodyParser(account); err != nil {
			return err
		}

		hashedPassword, err := hashPassword(account.Password)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		account.Password = hashedPassword

		account.SecretToken = generateSecretToken()

		// Map out all the claims to write to the payload
		claims := jwt.MapClaims{
			"id":  "TODO",
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}

		// Generate token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// TODO: Switch to RS256 with signing files https://github.com/gofiber/contrib/tree/main/jwt#rs256-example
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return &core.INTERNAL_SERVER_ERROR
		}

		if err := db.Create(&account).Error; err != nil {
			return &core.INTERNAL_SERVER_ERROR
		}

		// TODO: Implement refresh token
		return c.JSON(fiber.Map{"token": t})
	})
}

func hashPassword(password string) (string, error) {
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

func generateSecretToken() string {
	ll := len(ALLOWED_SECRET_TOKEN_CHARS)
	// 8 comes from db max length of secretToken
	result := make([]byte, SECRET_TOKEN_MAX_LENGTH)

	rand.Read(result)
	for i := 0; i < SECRET_TOKEN_MAX_LENGTH; i++ {
		result[i] = ALLOWED_SECRET_TOKEN_CHARS[int(result[i])%ll]
	}

	return string(result)
}
