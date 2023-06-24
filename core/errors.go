package core

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type RequestError struct {
	Code       int    `json:"code"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
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
)
