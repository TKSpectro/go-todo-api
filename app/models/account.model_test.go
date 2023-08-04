package models_test

import (
	"github.com/TKSpectro/go-todo-api/app/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Account.Model", func() {
	Describe("WriteRemote", func() {
		It("empty", func() {
			account := models.Account{}
			account.WriteRemote(models.Account{})

			Expect(account.Email).To(Equal(""))
			Expect(account.Password).To(Equal(""))
			Expect(account.Firstname).To(Equal(""))
			Expect(account.Lastname).To(Equal(""))
			Expect(account.TokenSecret).To(Equal(""))
		})

		It("with data", func() {
			account := models.Account{}
			account.WriteRemote(models.Account{
				Email:       "test@turbomeet.xyz",
				Password:    "password",
				Firstname:   "Firstname",
				Lastname:    "Lastname",
				TokenSecret: "token",
			})

			Expect(account.Email).To(Equal("test@turbomeet.xyz"))
			Expect(account.Password).To(Equal("")) // Password should not be written
			Expect(account.Firstname).To(Equal("Firstname"))
			Expect(account.Lastname).To(Equal("Lastname"))
			Expect(account.TokenSecret).To(Equal("")) // TokenSecret should not be written
		})
	})
})
