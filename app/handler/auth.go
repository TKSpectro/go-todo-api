package handler

import (
	"errors"
	"time"
	"tkspectro/vefeast/app/models"
	"tkspectro/vefeast/config/database"
	"tkspectro/vefeast/core"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	db := database.DB

	account := new(models.Account)
	if err := c.BodyParser(account); err != nil {
		return err
	}

	hashedPassword, err := models.HashPassword(account.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	account.Password = hashedPassword

	account.SecretToken = models.GenerateSecretToken()

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
}

func Login(c *fiber.Ctx) error {
	db := database.DB

	remote := new(models.Account)
	if err := c.BodyParser(remote); err != nil {
		return &core.BAD_REQUEST
	}

	var account models.Account
	if err := db.Where(&models.Account{Email: remote.Email}).Find(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return err
	}

	if !models.CheckPasswordHash(remote.Password, account.Password) {
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
}
