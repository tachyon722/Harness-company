package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockValidator struct {
	validateFn func(token string) (string, string, error)
}

func (m *mockValidator) ValidateToken(token string) (string, string, error) {
	return m.validateFn(token)
}

func TestRequireAuth_MissingHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := RequireAuth(&mockValidator{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRequireAuth_InvalidFormat(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mw := RequireAuth(&mockValidator{})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Basic token")
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRequireAuth_InvalidToken(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mw := RequireAuth(&mockValidator{
		validateFn: func(token string) (string, string, error) {
			return "", "", http.ErrAbortHandler
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestRequireAuth_Success(t *testing.T) {
	var capturedUserID, capturedUserType string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userInfo, ok := r.Context().Value(UserContextKey).(UserInfo)
		if !ok {
			t.Fatal("expected UserInfo in context")
		}
		capturedUserID = userInfo.UserID
		capturedUserType = userInfo.UserType
		w.WriteHeader(http.StatusOK)
	})

	mw := RequireAuth(&mockValidator{
		validateFn: func(token string) (string, string, error) {
			return "user-123", "human", nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	rec := httptest.NewRecorder()
	mw(handler).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if capturedUserID != "user-123" {
		t.Errorf("userID = %s, want user-123", capturedUserID)
	}
	if capturedUserType != "human" {
		t.Errorf("userType = %s, want human", capturedUserType)
	}
}
