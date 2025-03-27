package exc

import (
	"fmt"
	"net/http"
)

type Error struct {
	Err        string `json:"error"`
	Code       string `json:"code"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] - %s (%s)", e.StatusCode, e.Code, e.Message)
}

func NewAppError(err, code, message string, statusCode int) *Error {
	return &Error{
		Err:        err,
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
	}
}

func BadRequestError(code, msg string) *Error {
	return NewAppError("BadRequestError", code, msg, http.StatusBadRequest)
}

func NotFoundError(code, msg string) *Error {
	return NewAppError("NotFoundError", code, msg, http.StatusNotFound)
}

func InternalServerError(msg string) *Error {
	return NewAppError("InternalServerError", "internal_server_error", msg, http.StatusInternalServerError)
}

func UnauthorizedError(code, msg string) *Error {
	return NewAppError("UnauthorizedError", code, msg, http.StatusUnauthorized)
}

func ForbiddenError(code, msg string) *Error {
	return NewAppError("ForbiddenError", code, msg, http.StatusForbidden)
}

func ConflictError(code, msg string) *Error {
	return NewAppError("ConflictError", code, msg, http.StatusConflict)
}

func ServiceUnavailableError(code, msg string) *Error {
	return NewAppError("ServiceUnavailableError", code, msg, http.StatusServiceUnavailable)
}

func UnreachableOrigin(code, msg string) *Error {
	return NewAppError("UnreachableOrigin", code, msg, 523)
}

func RepositoryError(msg string) *Error {
	return InternalServerError(fmt.Sprintf("repository_error: [%s]", msg))
}

func ValidationError(code, field, msg string) *Error {
	return NewAppError("ValidationError", code, fmt.Sprintf("%s: [%s]", msg, field), http.StatusUnprocessableEntity)
}
