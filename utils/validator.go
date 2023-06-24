package utils

import (
	"fmt"
	"strings"
	"tkspectro/vefeast/core"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(payload interface{}) *core.RequestError {
	err := validate.Struct(payload)

	if err != nil {
		var errors []string

		for _, err := range err.(validator.ValidationErrors) {
			errors = append(
				errors,
				fmt.Sprintf("`%v` with value `%v` doesn't satisfy the `%v` constraint", err.Field(), err.Value(), err.Tag()),
			)
		}

		return &core.RequestError{
			Code:       core.VALIDATION_ERROR.Code,
			StatusCode: core.VALIDATION_ERROR.StatusCode,
			Message:    strings.Join(errors, ","),
		}
	}

	return nil
}
