package main

import (
	_ "github.com/TKSpectro/go-todo-api/api"
	"github.com/TKSpectro/go-todo-api/pkg/app"
	"github.com/TKSpectro/go-todo-api/pkg/database"
)

// @title           fiber-api
// @version         1.0
// @BasePath  /api
func main() {
	db := database.Connect()

	app := app.New(db)

	app.Listen(":3000")
}
