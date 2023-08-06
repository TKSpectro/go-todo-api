package model_test

import (
	"time"

	"github.com/TKSpectro/go-todo-api/app/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

var _ = Describe("Todo.Model", func() {
	Describe("WriteRemote", func() {
		It("empty", func() {
			todo := model.Todo{}
			todo.WriteRemote(model.Todo{})

			Expect(todo.Title).To(Equal(zero.NewString("", false)))
			Expect(todo.Description).To(Equal(zero.NewString("", false)))
			Expect(todo.Completed).To(Equal(false))
			Expect(todo.CompletedAt).To(Equal(null.NewTime(time.Time{}, false)))
		})

		It("with data", func() {
			time := time.Now()

			todo := model.Todo{}
			todo.WriteRemote(model.Todo{
				Title:       zero.NewString("Title", true),
				Description: zero.NewString("Description", true),
				Completed:   true,
				CompletedAt: null.NewTime(time, true),
			})

			Expect(todo.Title).To(Equal(zero.NewString("Title", true)))
			Expect(todo.Description).To(Equal(zero.NewString("Description", true)))
			Expect(todo.Completed).To(Equal(true))
			Expect(todo.CompletedAt).To(Equal(null.NewTime(time, true)))
		})
	})
})
