package app

import (
	_ "github.com/TKSpectro/go-todo-api/api"
	"github.com/TKSpectro/go-todo-api/config"
	"github.com/TKSpectro/go-todo-api/pkg/app/handler"
	"github.com/TKSpectro/go-todo-api/pkg/app/service"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html/v2"
)

func New(db *gorm.DB) *fiber.App {
	jwt.Init()

	engine := html.New(config.ROOT_PATH+"/pkg/view", ".html")

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
		Views:        engine,
		//! Can't use layouts in combination with rendering in create handlers (layout will get rendered too -> includes nav -> needs specific variables passed to it)
		//! This could maybe be solved by loading the template in the handler and executing the template manually. See: https://github.com/rngallen/gohtmx/blob/main/main.go
		// ViewsLayout:  "layouts/main",
		//  PassLocalsToViews: true, // TODO: Look into this
	})

	app.Use(logger.New())

	// Recover from panics anywhere in the chain and handle the control to the centralized ErrorHandler
	app.Use(recover.New())

	as := service.NewAccountService(db)
	ts := service.NewTodoService(db)

	h := handler.NewHandler(db, as, ts)

	h.RegisterRoutes(app)

	return app
}

func Shutdown(app *fiber.App) {
	app.Shutdown()
}
