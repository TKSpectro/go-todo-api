package app

import (
	_ "github.com/TKSpectro/go-todo-api/api"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/pkg/app/handler"
	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/service"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func New() *fiber.App {
	database.Connect()
	database.Migrate(&model.Account{}, &model.Todo{})

	jwt.Init()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use(logger.New())

	// Recover from panics anywhere in the chain and handle the control to the centralized ErrorHandler
	app.Use(recover.New())

	as := service.NewAccountService(database.DB)
	ts := service.NewTodoService(database.DB)

	h := handler.NewHandler(as, ts)

	h.RegisterRoutes(app)

	return app
}

func Shutdown(app *fiber.App) {
	database.Disconnect()

	app.Shutdown()
}
