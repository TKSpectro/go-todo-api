package handler

import (
	"fmt"
	"reflect"
	"strings"

	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"

	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := &Validator{
		validator: validator.New(),
	}

	v.RegisterCustomValidators()

	return v
}

func (v *Validator) Validate(payload interface{}) *utils.RequestError {
	if err := v.validator.Struct(payload); err != nil {
		var errors []string

		for _, err := range err.(validator.ValidationErrors) {
			errors = append(
				errors,
				fmt.Sprintf("`%v` with value `%v` doesn't satisfy the `%v` constraint", err.Field(), err.Value(), err.Tag()),
			)
		}

		return &utils.RequestError{
			Code:       utils.VALIDATION_ERROR.Code,
			StatusCode: utils.VALIDATION_ERROR.StatusCode,
			Message:    strings.Join(errors, ","),
		}
	}

	return nil
}

func (v *Validator) RegisterCustomValidators() {

	// Register custom validators for zero and null types
	v.validator.RegisterCustomTypeFunc(ValidateZeroString, zero.String{})
	v.validator.RegisterCustomTypeFunc(ValidateZeroInt, zero.Int{})
	v.validator.RegisterCustomTypeFunc(ValidateZeroBool, zero.Bool{})
	v.validator.RegisterCustomTypeFunc(ValidateZeroFloat, zero.Float{})
	v.validator.RegisterCustomTypeFunc(ValidateZeroTime, zero.Time{})
	v.validator.RegisterCustomTypeFunc(ValidateNullString, null.String{})
	v.validator.RegisterCustomTypeFunc(ValidateNullInt, null.Int{})
	v.validator.RegisterCustomTypeFunc(ValidateNullBool, null.Bool{})
	v.validator.RegisterCustomTypeFunc(ValidateNullFloat, null.Float{})
	v.validator.RegisterCustomTypeFunc(ValidateNullTime, null.Time{})
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
