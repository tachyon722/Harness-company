package errors

import (
	"net/http"
	"testing"
)

func TestAppError_HTTPStatus(t *testing.T) {
	tests := []struct {
		name     string
		code     Code
		wantCode int
	}{
		{"invalid request", ErrInvalidRequest, http.StatusBadRequest},
		{"validation", ErrValidation, http.StatusBadRequest},
		{"unauthorized", ErrUnauthorized, http.StatusUnauthorized},
		{"invalid token", ErrInvalidToken, http.StatusUnauthorized},
		{"token expired", ErrTokenExpired, http.StatusUnauthorized},
		{"forbidden", ErrForbidden, http.StatusForbidden},
		{"not found", ErrNotFound, http.StatusNotFound},
		{"conflict", ErrConflict, http.StatusConflict},
		{"internal", ErrInternal, http.StatusInternalServerError},
		{"unknown", "UNKNOWN", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &AppError{Code: tt.code, Message: "test"}
			if got := err.HTTPStatus(); got != tt.wantCode {
				t.Errorf("HTTPStatus() = %v, want %v", got, tt.wantCode)
			}
		})
	}
}

func TestAppError_Error(t *testing.T) {
	msg := "something went wrong"
	err := &AppError{Code: ErrInternal, Message: msg}
	if got := err.Error(); got != msg {
		t.Errorf("Error() = %v, want %v", got, msg)
	}
}

func TestNewFunctions(t *testing.T) {
	tests := []struct {
		name string
		got  *AppError
		code Code
	}{
		{"NewInvalidRequest", NewInvalidRequest("bad"), ErrInvalidRequest},
		{"NewUnauthorized", NewUnauthorized("no access"), ErrUnauthorized},
		{"NewForbidden", NewForbidden("blocked"), ErrForbidden},
		{"NewNotFound", NewNotFound("missing"), ErrNotFound},
		{"NewConflict", NewConflict("duplicate"), ErrConflict},
		{"NewValidation", NewValidation("invalid"), ErrValidation},
		{"NewInternal", NewInternal("oops"), ErrInternal},
		{"NewTokenExpired", NewTokenExpired(), ErrTokenExpired},
		{"NewInvalidToken", NewInvalidToken(), ErrInvalidToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got.Code != tt.code {
				t.Errorf("Code = %v, want %v", tt.got.Code, tt.code)
			}
		})
	}
}
