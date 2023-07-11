package core

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RequestError struct {
	Code       int    `json:"code"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Detail     string `json:"detail,omitempty"`
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("statusCode: %d | code: %d | message: %v", r.StatusCode, r.Code, r.Message)
}

var (
	BAD_REQUEST           = RequestError{Code: 400, StatusCode: fiber.StatusBadRequest, Message: "Missing params/body."}
	NOT_FOUND             = RequestError{Code: 404, StatusCode: fiber.StatusNotFound, Message: "Not found."}
	UNAUTHORIZED          = RequestError{Code: 401, StatusCode: fiber.StatusUnauthorized, Message: "Unauthorized."}
	FORBIDDEN             = RequestError{Code: 403, StatusCode: fiber.StatusForbidden, Message: "Forbidden."}
	NOT_ACCEPTABLE        = RequestError{Code: 406, StatusCode: fiber.StatusNotAcceptable, Message: "Not acceptable."}
	INTERNAL_SERVER_ERROR = RequestError{Code: 500, StatusCode: fiber.StatusInternalServerError, Message: "Internal Server error."}

	// The validation error can be used but the message should be overwritten
	VALIDATION_ERROR = RequestError{Code: 1000, StatusCode: fiber.StatusBadRequest, Message: "Validation error."}

	AUTH_LOGIN_WRONG_PASSWORD = RequestError{Code: 1050, StatusCode: fiber.StatusUnauthorized, Message: "Wrong password."}
	WRONG_REFRESH_TOKEN       = RequestError{Code: 1051, StatusCode: fiber.StatusUnauthorized, Message: "Wrong refresh token."}

	ACCOUNT_WITH_EMAIL_ALREADY_EXISTS = RequestError{Code: 1100, StatusCode: fiber.StatusBadRequest, Message: "An account with this email already exists."}
)

// Error from var Error but pass details
func RequestErrorFrom(err *RequestError, detail string) *RequestError {
	err.Detail = detail
	return err
}
