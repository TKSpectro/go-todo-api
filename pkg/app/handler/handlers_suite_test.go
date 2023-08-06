package handler_test

import (
	"testing"

	"github.com/TKSpectro/go-todo-api/pkg/app"
	"github.com/TKSpectro/go-todo-api/pkg/database"
	"github.com/TKSpectro/go-todo-api/test"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

var App *fiber.App
var DB *gorm.DB

var _ = BeforeSuite(func() {
	test.New()

	DB = database.Connect()
	App = app.New(DB)
})

var _ = AfterSuite(func() {
	test.ClearAllTables()

	app.Shutdown(App)
	database.Disconnect(DB)
})
