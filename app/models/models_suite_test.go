package models_test

import (
	"testing"

	"github.com/TKSpectro/go-todo-api/pkg/app"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModels(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var testApp *fiber.App

var _ = BeforeSuite(func() {
	testApp = app.Setup()
})

var _ = AfterSuite(func() {
	app.Shutdown(testApp)
})
