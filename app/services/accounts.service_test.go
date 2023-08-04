package services_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Accounts.Service", func() {
	Describe("List", func() {
		It("empty", func() {
			req, _ := http.NewRequest("GET", "/api/accounts", nil)
			res, _ := App.Test(req)

			Expect(res.StatusCode).To(Equal(401))
		})
	})
})
