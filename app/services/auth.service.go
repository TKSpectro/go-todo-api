package services

import (
	"errors"
	"tkspectro/vefeast/app/models"
	"tkspectro/vefeast/app/types"
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
		SecretToken: models.GenerateSecretToken(),
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
		ID: account.ID,
	})

	refreshToken := "TODO"

	// TODO: Implement refresh token
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

	refreshToken := "TODO"

	// TODO: Implement refresh token
	return c.JSON(&types.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}
