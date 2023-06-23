package router

import (
	"crypto/rand"
	"time"
	"tkspectro/vefeast/core"
	"tkspectro/vefeast/model"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterRoutesAuth(router fiber.Router) {
	// Authentication routes
	router.Put("/login", func(c *fiber.Ctx) error {
		remote := new(LoginBody)

		if err := c.BodyParser(remote); err != nil {
			return err
		}

		// TODO: Get user from db
		if false {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

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
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// TODO: Implement refresh token
		return c.JSON(fiber.Map{"token": t})
	})

	router.Post("/register", func(c *fiber.Ctx) error {
		db := core.DB

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
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := db.Create(&account).Error; err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// TODO: Implement refresh token
		return c.JSON(fiber.Map{"token": t})
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
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

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
