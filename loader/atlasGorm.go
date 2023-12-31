package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/TKSpectro/go-todo-api/pkg/app/model"
)

func main() {
	stmts, err := gormschema.New("mysql").Load(&model.Account{}, &model.Todo{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
