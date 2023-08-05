package services_test

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("Accounts.Service", func() {
	Describe("List", Ordered, func() {
		var authToken string
		BeforeAll(func() {
			pw, _ := models.HashPassword("123456")
			account := models.Account{
				Email:     "accounts.list@turbomeet.xyz",
				Password:  pw,
				Firstname: "Accounts",
				Lastname:  "List",
			}
			models.CreateAccount(&account)

			auth, _ := jwt.Generate(&account)
			authToken = auth.Token
		})

		It("should be unauthorized", func() {
			req, _ := http.NewRequest("GET", "/api/accounts", nil)
			res, _ := App.Test(req)

			Expect(res.StatusCode).To(Equal(401))
		})

		It("should be authorized", func() {
			req, _ := http.NewRequest("GET", "/api/accounts", nil)
			req.Header.Set("Authorization", "Bearer "+authToken)
			res, _ := App.Test(req)

			Expect(res.StatusCode).To(Equal(200))

			bodyBytes, _ := io.ReadAll(res.Body)
			result := types.GetAccountsResponse{}
			if err := json.Unmarshal(bodyBytes, &result); err != nil {
				PanicWith("Error unmarshalling response")
			}

			Expect(result.Accounts).To(ContainElement(MatchFields(IgnoreExtras, Fields{
				"Email":     Equal("accounts.list@turbomeet.xyz"),
				"Password":  Equal(""), // Password should be empty in response
				"Firstname": Equal("Accounts"),
				"Lastname":  Equal("List"),
			})))
		})

		AfterAll(func() {
			database.DB.Exec("DELETE FROM accounts")
		})
	})
})
