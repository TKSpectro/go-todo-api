package services

import (
	"errors"

	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/TKSpectro/go-todo-api/utils/jwt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Register      godoc
// @Summary      Register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        account body types.RegisterDTO true "Account"
// @Success      200 {object} types.AuthResponse
// @Router       /auth/register [post]
func Register(c *fiber.Ctx) error {
	remote := new(types.RegisterDTO)

	if err := utils.ParseBodyAndValidate(c, remote); err != nil {
		return err
	}

	err := models.FindAccountByEmail(&struct{ ID string }{}, remote.Account.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return &core.ACCOUNT_WITH_EMAIL_ALREADY_EXISTS
	}

	account := &models.Account{
		Email:       remote.Account.Email,
		TokenSecret: models.GenerateSecretToken(),
		Firstname:   remote.Account.Firstname,
		Lastname:    remote.Account.Lastname,
	}

	hashedPassword, err := models.HashPassword(remote.Account.Password)
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
		Auth: types.AuthResponseBody{
			Token:        token,
			RefreshToken: refreshToken,
		},
	})
}

// Login      godoc
// @Summary      Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        account body types.LoginDTO true "Account"
// @Success      200 {object} types.AuthResponse
// @Router       /auth/login [put]
func Login(c *fiber.Ctx) error {
	remote := new(types.LoginDTO)

	if err := utils.ParseBodyAndValidate(c, remote); err != nil {
		return err
	}

	account := &models.Account{}
	if err := models.FindAccountByEmail(account, remote.Account.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return err
	}

	if !models.CheckPasswordHash(remote.Account.Password, account.Password) {
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
		Auth: types.AuthResponseBody{
			Token:        token,
			RefreshToken: refreshToken,
		},
	})
}

// Refresh      godoc
// @Summary      Refresh
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200 {object} types.AuthResponse
// @Router       /auth/refresh [put]
func Refresh(c *fiber.Ctx) error {
	var accountId = c.Locals("AccountId").(uint)
	var tokenType = c.Locals("TokenType").(string)
	var tokenSecret = c.Locals("TokenSecret").(string)

	if tokenType != "refresh" {
		return &core.WRONG_REFRESH_TOKEN
	}

	account := &models.Account{}
	if err := database.DB.Model(account).Take(account, &models.Account{
		BaseModel:   models.BaseModel{ID: accountId},
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
		Auth: types.AuthResponseBody{
			Token:        token,
			RefreshToken: refreshToken,
		},
	})
}
