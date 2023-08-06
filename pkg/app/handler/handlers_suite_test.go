package handler_test

import (
	"testing"

	"github.com/TKSpectro/go-todo-api/pkg/app"
	"github.com/TKSpectro/go-todo-api/test"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

var App *fiber.App

var _ = BeforeSuite(func() {
	test.New()

	App = app.New()
})

var _ = AfterSuite(func() {
	test.ClearAllTables()

	app.Shutdown(App)
})