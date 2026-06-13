package errors

import "net/http"

type Code string

const (
	ErrInvalidRequest   Code = "INVALID_REQUEST"
	ErrUnauthorized     Code = "UNAUTHORIZED"
	ErrForbidden        Code = "FORBIDDEN"
	ErrNotFound         Code = "NOT_FOUND"
	ErrConflict         Code = "CONFLICT"
	ErrValidation       Code = "VALIDATION_ERROR"
	ErrInternal         Code = "INTERNAL_ERROR"
	ErrTokenExpired     Code = "TOKEN_EXPIRED"
	ErrInvalidToken     Code = "INVALID_TOKEN"
)

type AppError struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrInvalidRequest, ErrValidation:
		return http.StatusBadRequest
	case ErrUnauthorized, ErrInvalidToken, ErrTokenExpired:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func NewInvalidRequest(msg string) *AppError {
	return &AppError{Code: ErrInvalidRequest, Message: msg}
}

func NewUnauthorized(msg string) *AppError {
	return &AppError{Code: ErrUnauthorized, Message: msg}
}

func NewForbidden(msg string) *AppError {
	return &AppError{Code: ErrForbidden, Message: msg}
}

func NewNotFound(msg string) *AppError {
	return &AppError{Code: ErrNotFound, Message: msg}
}

func NewConflict(msg string) *AppError {
	return &AppError{Code: ErrConflict, Message: msg}
}

func NewValidation(msg string) *AppError {
	return &AppError{Code: ErrValidation, Message: msg}
}

func NewInternal(msg string) *AppError {
	return &AppError{Code: ErrInternal, Message: msg}
}

func NewTokenExpired() *AppError {
	return &AppError{Code: ErrTokenExpired, Message: "token has expired"}
}

func NewInvalidToken() *AppError {
	return &AppError{Code: ErrInvalidToken, Message: "invalid token"}
}
