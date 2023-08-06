package model_test

import (
	"testing"

	"github.com/TKSpectro/go-todo-api/pkg/app"
	"github.com/TKSpectro/go-todo-api/test"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model Suite")
}

var testApp *fiber.App

var _ = BeforeSuite(func() {
	test.New()

	testApp = app.New()
})

var _ = AfterSuite(func() {
	test.ClearAllTables()

	app.Shutdown(testApp)
})