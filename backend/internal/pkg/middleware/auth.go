package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	apperrors "github.com/harness-org/backend/internal/pkg/errors"
)

type contextKey string

const UserContextKey contextKey = "user"

type UserInfo struct {
	UserID   string `json:"user_id"`
	UserType string `json:"user_type"`
}

type TokenValidator interface {
	ValidateToken(tokenString string) (string, string, error)
}

func RequireAuth(validator TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, apperrors.NewUnauthorized("missing authorization header"))
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token == authHeader {
				writeError(w, apperrors.NewUnauthorized("invalid authorization format"))
				return
			}

			userID, userType, err := validator.ValidateToken(token)
			if err != nil {
				errMsg := err.Error()
				if strings.Contains(errMsg, "expired") {
					writeError(w, apperrors.NewTokenExpired())
				} else {
					writeError(w, apperrors.NewInvalidToken())
				}
				return
			}

			userInfo := UserInfo{UserID: userID, UserType: userType}
			ctx := context.WithValue(r.Context(), UserContextKey, userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeError(w http.ResponseWriter, appErr *apperrors.AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus())
	json.NewEncoder(w).Encode(map[string]any{
		"error": appErr.Message,
		"code":  appErr.Code,
	})
}
