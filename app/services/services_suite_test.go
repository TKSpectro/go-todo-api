package services_test

import (
	"testing"

	"github.com/TKSpectro/go-todo-api/pkg/app"
	"github.com/TKSpectro/go-todo-api/test"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

var App *fiber.App

var _ = BeforeSuite(func() {
	test.Setup()

	App = app.Setup()
})

var _ = AfterSuite(func() {
	test.ClearAllTables()

	app.Shutdown(App)
})
