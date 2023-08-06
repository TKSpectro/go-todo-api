package app

import (
	_ "github.com/TKSpectro/go-todo-api/api"
	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/routes"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func New() *fiber.App {
	database.Connect()
	database.Migrate(&models.Account{}, &models.Todo{})

	jwt.Init()

	utils.RegisterCustomValidators()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use(logger.New())

	// Recover from panics anywhere in the chain and handle the control to the centralized ErrorHandler
	app.Use(recover.New())

	routes.New(app)

	return app
}

func Shutdown(app *fiber.App) {
	database.Disconnect()

	app.Shutdown()
}
