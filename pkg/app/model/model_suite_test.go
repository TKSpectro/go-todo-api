package model_test

import (
	"testing"

	"github.com/TKSpectro/go-todo-api/pkg/app"
	"github.com/TKSpectro/go-todo-api/test"
	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

func TestModel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Model Suite")
}

var testApp *fiber.App
var DB *gorm.DB

var _ = BeforeSuite(func() {
	DB = test.Setup()
	testApp = app.New(DB)
})

var _ = AfterSuite(func() {
	app.Shutdown(testApp)
	test.Teardown(DB)
})
