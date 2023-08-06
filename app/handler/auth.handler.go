package handler

import (
	"errors"

	"github.com/TKSpectro/go-todo-api/app/model"
	"github.com/TKSpectro/go-todo-api/app/types"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/utils"
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
	remoteData := &types.RegisterDTO{}

	if err := utils.ParseBodyAndValidate(c, remoteData); err != nil {
		return err
	}

	err := model.FindAccountByEmail(struct{}{}, remoteData.Account.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return &core.ACCOUNT_WITH_EMAIL_ALREADY_EXISTS
	}

	account := &model.Account{}
	// Convert remoteData.Account from RegisterDTOBody to Account type
	account.WriteRemote(utils.Convert(model.Account{}, &remoteData.Account))

	account.TokenSecret = model.GenerateSecretToken()
	hashedPassword, err := model.HashPassword(remoteData.Account.Password)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	account.Password = hashedPassword

	if err := model.CreateAccount(account).Error; err != nil {
		return &core.INTERNAL_SERVER_ERROR
	}

	auth, err := jwt.Generate(account)
	if err != nil {
		return err
	}

	return c.JSON(&types.AuthResponse{
		Auth: auth,
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
	remoteData := &types.LoginDTO{}

	if err := utils.ParseBodyAndValidate(c, remoteData); err != nil {
		return err
	}

	account := &model.Account{}
	if err := model.FindAccountByEmail(account, remoteData.Account.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return err
	}

	if !model.CheckPasswordHash(remoteData.Account.Password, account.Password) {
		return &core.AUTH_LOGIN_WRONG_PASSWORD
	}

	auth, err := jwt.Generate(account)
	if err != nil {
		return err
	}

	return c.JSON(&types.AuthResponse{
		Auth: auth,
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
	var tokenPayload = locals.JwtPayload(c)

	if tokenPayload.Type != "refresh" {
		return &core.WRONG_REFRESH_TOKEN
	}

	account := &model.Account{}
	if err := database.DB.Model(account).Take(account, &model.Account{
		BaseModel:   model.BaseModel{ID: tokenPayload.AccountID},
		TokenSecret: tokenPayload.Secret,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.UNAUTHORIZED
		}
		return err
	}

	auth, err := jwt.Generate(account)
	if err != nil {
		return err
	}

	return c.JSON(&types.AuthResponse{
		Auth: auth,
	})
}

func RotateJWK(c *fiber.Ctx) error {
	jwt.RotateJWK()

	return c.SendStatus(fiber.StatusOK)
}

func Me(c *fiber.Ctx) error {
	var account = &model.Account{}

	err := model.FindAccountByID(account, locals.JwtPayload(c).AccountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetAccountResponse{
		Account: *account,
	})
}
