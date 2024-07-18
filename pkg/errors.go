package pkg

import (
	"fmt"
	"net/http"
)

type ErrorCode int

const (
	NotFound = 1000 + iota
	BadRequest
	InternalError
	WrongFieldValues
	Unauthorized
	Forbidden
)

type ApiError struct {
	Code         ErrorCode
	ResponseCode int
	Message      string
}

func (e ApiError) Error() string {
	return fmt.Sprintf(`{"error_code": %d, "message": "%s"}`, e.Code, e.Message)
}

func GetError(code ErrorCode) ApiError {
	return apiErrors[code]
}

var apiErrors = map[ErrorCode]ApiError{
	NotFound: {
		Code:         NotFound,
		Message:      "Resource not found",
		ResponseCode: http.StatusNotFound,
	},
	BadRequest: {
		Code:         BadRequest,
		Message:      "Bad request",
		ResponseCode: http.StatusBadRequest,
	},
	InternalError: {
		Code:         InternalError,
		Message:      "Internal server error",
		ResponseCode: http.StatusInternalServerError,
	},
	WrongFieldValues: {
		Code:         WrongFieldValues,
		Message:      "invalid field values",
		ResponseCode: http.StatusBadRequest,
	},
	Unauthorized: {
		Code:         Unauthorized,
		Message:      "Unauthorized",
		ResponseCode: http.StatusUnauthorized,
	},
	Forbidden: {
		Code:         Forbidden,
		Message:      "Forbidden",
		ResponseCode: http.StatusForbidden,
	},
}
