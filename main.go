package main

import (
	_ "github.com/TKSpectro/go-todo-api/api"
	"github.com/TKSpectro/go-todo-api/pkg/app"
)

// @title           fiber-api
// @version         1.0
// @BasePath  /api
func main() {
	app := app.New()

	app.Listen(":3000")
}
