package services

import (
	"errors"
	"tkspectro/vefeast/app/models"
	"tkspectro/vefeast/app/types"
	"tkspectro/vefeast/config/database"
	"tkspectro/vefeast/core"
	"tkspectro/vefeast/utils"
	"tkspectro/vefeast/utils/jwt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	remote := new(types.RegisterDTO)

	if err := utils.ParseBodyAndValidate(c, remote); err != nil {
		return err
	}

	err := models.FindAccountByEmail(&struct{ ID string }{}, remote.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return &core.ACCOUNT_WITH_EMAIL_ALREADY_EXISTS
	}

	account := &models.Account{
		Email:       remote.Email,
		TokenSecret: models.GenerateSecretToken(),
		Firstname:   remote.Firstname,
		Lastname:    remote.Lastname,
	}

	hashedPassword, err := models.HashPassword(remote.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	account.Password = hashedPassword

	if err := models.CreateAccount(account).Error; err != nil {
		return &core.INTERNAL_SERVER_ERROR
	}

	token := jwt.Generate(&jwt.TokenPayload{
		ID:   account.ID,
		Type: "auth",
	})

	refreshToken := jwt.Generate(&jwt.TokenPayload{
		ID:     account.ID,
		Type:   "refresh",
		Secret: account.TokenSecret,
	})

	return c.JSON(&types.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func Login(c *fiber.Ctx) error {
	remote := new(types.LoginDTO)

	if err := utils.ParseBodyAndValidate(c, remote); err != nil {
		return err
	}

	account := &models.Account{}
	if err := models.FindAccountByEmail(account, remote.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return err
	}

	if !models.CheckPasswordHash(remote.Password, account.Password) {
		return &core.AUTH_LOGIN_WRONG_PASSWORD
	}

	token := jwt.Generate(&jwt.TokenPayload{
		ID: account.ID,
	})

	refreshToken := jwt.Generate(&jwt.TokenPayload{
		ID:     account.ID,
		Type:   "refresh",
		Secret: account.TokenSecret,
	})

	return c.JSON(&types.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func Refresh(c *fiber.Ctx) error {
	var accountId = c.Locals("AccountId").(uint)
	var tokenType = c.Locals("TokenType").(string)
	var tokenSecret = c.Locals("TokenSecret").(string)

	if tokenType != "refresh" {
		return &core.WRONG_REFRESH_TOKEN
	}

	// TODO: Find account by id and tokenSecret

	account := &models.Account{}
	if err := database.DB.Model(account).Take(account, &models.Account{
		Model:       gorm.Model{ID: accountId},
		TokenSecret: tokenSecret,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.UNAUTHORIZED
		}
		return err
	}

	token := jwt.Generate(&jwt.TokenPayload{
		ID: account.ID,
	})

	refreshToken := jwt.Generate(&jwt.TokenPayload{
		ID:     account.ID,
		Type:   "refresh",
		Secret: account.TokenSecret,
	})

	return c.JSON(&types.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}
