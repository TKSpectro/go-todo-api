package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/TKSpectro/go-todo-api/core"

	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/go-playground/validator/v10"
)

var Validator = validator.New()

func Validate(payload interface{}) *core.RequestError {
	err := Validator.Struct(payload)

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

func RegisterCustomValidators() {

	// Register custom validators for zero and null types
	Validator.RegisterCustomTypeFunc(ValidateZeroString, zero.String{})
	Validator.RegisterCustomTypeFunc(ValidateZeroInt, zero.Int{})
	Validator.RegisterCustomTypeFunc(ValidateZeroBool, zero.Bool{})
	Validator.RegisterCustomTypeFunc(ValidateZeroFloat, zero.Float{})
	Validator.RegisterCustomTypeFunc(ValidateZeroTime, zero.Time{})
	Validator.RegisterCustomTypeFunc(ValidateNullString, null.String{})
	Validator.RegisterCustomTypeFunc(ValidateNullInt, null.Int{})
	Validator.RegisterCustomTypeFunc(ValidateNullBool, null.Bool{})
	Validator.RegisterCustomTypeFunc(ValidateNullFloat, null.Float{})
	Validator.RegisterCustomTypeFunc(ValidateNullTime, null.Time{})
}

func ValidateZeroString(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(zero.String); ok {
		if valuer.Valid {
			return valuer.String
		}
	}

	return nil
}

func ValidateZeroInt(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(zero.Int); ok {
		if valuer.Valid {
			return valuer.Int64
		}
	}

	return nil
}

func ValidateZeroBool(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(zero.Bool); ok {
		if valuer.Valid {
			return valuer.Bool
		}
	}

	return nil
}

func ValidateZeroFloat(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(zero.Float); ok {
		if valuer.Valid {
			return valuer.Float64
		}
	}

	return nil
}

func ValidateZeroTime(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(zero.Time); ok {
		if valuer.Valid {
			return valuer.Time
		}
	}

	return nil
}

func ValidateNullString(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.String); ok {
		if valuer.Valid {
			return valuer.String
		}
	}

	return nil
}

func ValidateNullInt(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.Int); ok {
		if valuer.Valid {
			return valuer.Int64
		}
	}

	return nil
}

func ValidateNullBool(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.Bool); ok {
		if valuer.Valid {
			return valuer.Bool
		}
	}

	return nil
}

func ValidateNullFloat(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.Float); ok {
		if valuer.Valid {
			return valuer.Float64
		}
	}

	return nil
}

func ValidateNullTime(field reflect.Value) interface{} {
	if valuer, ok := field.Interface().(null.Time); ok {
		if valuer.Valid {
			return valuer.Time
		}
	}

	return nil
}
