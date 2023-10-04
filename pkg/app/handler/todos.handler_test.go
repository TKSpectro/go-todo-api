package handler_test

import (
	"net/http"

	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/service"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// . "github.com/onsi/gomega/gstruct"
	"gopkg.in/guregu/null.v4/zero"
)

var _ = Describe("Todos.Handler", Ordered, func() {
	Describe("ExportCSV", Ordered, func() {
		var authToken string
		BeforeAll(func() {
			pw, _ := model.HashPassword("123456")
			account := &model.Account{
				Email:     "todos.exportcsv@turbomeet.xyz",
				Password:  pw,
				Firstname: "Accounts",
				Lastname:  "ExportCSV",
			}
			accountService := service.NewAccountService(DB)
			accountService.CreateAccount(account)

			auth, _ := jwt.Generate(account)
			authToken = auth.Token

			todoService := service.NewTodoService(DB)
			todoService.CreateTodo(&model.Todo{
				Title:     zero.NewString("Todo 1", true),
				Completed: false,
				AccountID: account.ID,
			})

			todoService.CreateTodo(&model.Todo{
				Title:       zero.NewString("Todo 2", true),
				Completed:   false,
				AccountID:   account.ID,
				Description: zero.NewString("some description", true),
			})

			todoService.CreateTodo(&model.Todo{
				Title:       zero.NewString("Todo 3", true),
				Completed:   true,
				AccountID:   account.ID,
				Description: zero.NewString("some longer description with a \n line-break", true),
			})

			todoService.CreateTodo(&model.Todo{
				Title:     zero.NewString("Todo 4", true),
				Completed: true,
				AccountID: account.ID,
			})
		})

		It("should be unauthorized", func() {
			req, _ := http.NewRequest("GET", "/api/todos/csv", nil)
			res, _ := App.Test(req)

			Expect(res.StatusCode).To(Equal(401))
		})

		It("should be authorized", func() {
			req, _ := http.NewRequest("GET", "/api/todos/csv", nil)
			req.Header.Set("Authorization", "Bearer "+authToken)
			res, _ := App.Test(req)

			Expect(res.StatusCode).To(Equal(200))
			Expect(res.Header.Get("Content-Type")).To(Equal("text/csv"))

			// With the following code you could write the response to a file
			// fo, err := os.Create(config.TEST_FILE_PATH + "todos-export.csv")
			// if err != nil {
			// 	panic(err)
			// }

			// // close fo on exit and check for its returned error
			// defer func() {
			// 	if err := fo.Close(); err != nil {
			// 		panic(err)
			// 	}
			// }()

			// // make a buffer to keep chunks that are read
			// buf := make([]byte, 1024)
			// for {
			// 	// read a chunk
			// 	n, err := res.Body.Read(buf)
			// 	if err != nil && err != io.EOF {
			// 		panic(err)
			// 	}
			// 	if n == 0 {
			// 		break
			// 	}

			// 	// write a chunk
			// 	if _, err := fo.Write(buf[:n]); err != nil {
			// 		panic(err)
			// 	}
			// }
		})

		AfterAll(func() {
			test.ClearTables(DB, []string{"accounts"})
		})
	})

	AfterAll(func() {
		test.ClearAllTables(DB)
	})
})
