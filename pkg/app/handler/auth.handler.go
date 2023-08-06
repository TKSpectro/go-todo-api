package handler

import (
	"errors"

	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types"
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
func (h *Handler) Register(c *fiber.Ctx) error {
	remoteData := &types.RegisterDTO{}

	if err := ParseBodyAndValidate(c, remoteData, *h.validator); err != nil {
		return err
	}

	err := h.accountService.FindAccountByEmail(struct{}{}, remoteData.Account.Email).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return &utils.ACCOUNT_WITH_EMAIL_ALREADY_EXISTS
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

	if err := h.accountService.CreateAccount(account).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
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
func (h *Handler) Login(c *fiber.Ctx) error {
	remoteData := &types.LoginDTO{}

	if err := ParseBodyAndValidate(c, remoteData, *h.validator); err != nil {
		return err
	}

	account := &model.Account{}
	if err := h.accountService.FindAccountByEmail(account, remoteData.Account.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return err
	}

	if !model.CheckPasswordHash(remoteData.Account.Password, account.Password) {
		return &utils.AUTH_LOGIN_WRONG_PASSWORD
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
func (h *Handler) Refresh(c *fiber.Ctx) error {
	var tokenPayload = locals.JwtPayload(c)

	if tokenPayload.Type != "refresh" {
		return &utils.WRONG_REFRESH_TOKEN
	}

	account := &model.Account{}
	if err := h.db.Model(account).Take(account, &model.Account{
		BaseModel:   model.BaseModel{ID: tokenPayload.AccountID},
		TokenSecret: tokenPayload.Secret,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.UNAUTHORIZED
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

func (h *Handler) RotateJWK(c *fiber.Ctx) error {
	jwt.RotateJWK()

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Me(c *fiber.Ctx) error {
	var account = &model.Account{}

	err := h.accountService.FindAccountByID(account, locals.JwtPayload(c).AccountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetAccountResponse{
		Account: *account,
	})
}
