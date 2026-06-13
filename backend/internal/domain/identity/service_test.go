package identity

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type mockRepo struct {
	createUserFn    func(ctx context.Context, input CreateUserInput) (*User, error)
	getUserByEmailFn func(ctx context.Context, email string) (*User, error)
	getUserByIDFn   func(ctx context.Context, id uuid.UUID) (*User, error)
	createAgentFn   func(ctx context.Context, input CreateAgentInput) (*AIAgent, string, error)
	getAgentByIDFn  func(ctx context.Context, id uuid.UUID) (*AIAgent, error)
	listRolesFn     func(ctx context.Context) ([]Role, error)
}

func (m *mockRepo) CreateUser(ctx context.Context, input CreateUserInput) (*User, error) {
	return m.createUserFn(ctx, input)
}
func (m *mockRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return m.getUserByEmailFn(ctx, email)
}
func (m *mockRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return m.getUserByIDFn(ctx, id)
}
func (m *mockRepo) CreateAgent(ctx context.Context, input CreateAgentInput) (*AIAgent, string, error) {
	return m.createAgentFn(ctx, input)
}
func (m *mockRepo) GetAgentByID(ctx context.Context, id uuid.UUID) (*AIAgent, error) {
	return m.getAgentByIDFn(ctx, id)
}
func (m *mockRepo) ListRoles(ctx context.Context) ([]Role, error) {
	return m.listRolesFn(ctx)
}

func TestRegisterUser_Validation(t *testing.T) {
	svc := NewService(&mockRepo{}, "secret")
	tests := []struct {
		name  string
		input CreateUserInput
	}{
		{"empty name", CreateUserInput{Email: "a@b.com", Password: "pass"}},
		{"empty email", CreateUserInput{Name: "a", Password: "pass"}},
		{"empty password", CreateUserInput{Name: "a", Email: "a@b.com"}},
		{"all empty", CreateUserInput{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.RegisterUser(context.Background(), tt.input)
			if !errors.Is(err, ErrValidation) {
				t.Errorf("expected ErrValidation, got %v", err)
			}
		})
	}
}

func TestRegisterUser_Success(t *testing.T) {
	uid := uuid.New()
	now := time.Now()
	svc := NewService(&mockRepo{
		createUserFn: func(ctx context.Context, input CreateUserInput) (*User, error) {
			return &User{ID: uid, Name: input.Name, Email: input.Email, CreatedAt: now, UpdatedAt: now}, nil
		},
	}, "secret")

	resp, err := svc.RegisterUser(context.Background(), CreateUserInput{Name: "Alice", Email: "alice@test.com", Password: "secure123"})
	if err != nil {
		t.Fatalf("RegisterUser() error = %v", err)
	}
	if resp.Name != "Alice" || resp.Email != "alice@test.com" {
		t.Errorf("got %+v, want Name=Alice Email=alice@test.com", resp)
	}
}

func TestAuthenticateUser_InvalidCredentials(t *testing.T) {
	svc := NewService(&mockRepo{
		getUserByEmailFn: func(ctx context.Context, email string) (*User, error) {
			return nil, errors.New("not found")
		},
	}, "secret")

	_, err := svc.AuthenticateUser(context.Background(), "nonexist@test.com", "password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthenticateUser_EmptyInput(t *testing.T) {
	svc := NewService(&mockRepo{}, "secret")

	_, err := svc.AuthenticateUser(context.Background(), "", "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestAuthenticateAgent_EmptyKey(t *testing.T) {
	svc := NewService(&mockRepo{}, "secret")

	_, err := svc.AuthenticateAgent(context.Background(), uuid.New(), "")
	if !errors.Is(err, ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc := NewService(&mockRepo{}, "secret")

	_, _, err := svc.ValidateToken("invalid.token.string")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	svc := NewService(&mockRepo{}, "my-secret-key")
	svc.tokenTTL = time.Hour

	token, expiresAt, err := svc.generateJWT("user-123", "human", "TestUser")
	if err != nil {
		t.Fatalf("generateJWT() error = %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
	if expiresAt <= time.Now().Unix() {
		t.Fatal("expected future expiration")
	}

	userID, userType, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if userID != "user-123" {
		t.Errorf("userID = %v, want user-123", userID)
	}
	if userType != "human" {
		t.Errorf("userType = %v, want human", userType)
	}
}

func TestTokenExpired(t *testing.T) {
	svc := NewService(&mockRepo{}, "secret")
	svc.tokenTTL = -time.Hour

	token, _, err := svc.generateJWT("user-1", "human", "Test")
	if err != nil {
		t.Fatalf("generateJWT() error = %v", err)
	}

	_, _, err = svc.ValidateToken(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestListRoles(t *testing.T) {
	expectedRoles := []Role{
		{Name: "admin", RoleType: RolePlanner},
		{Name: "exec", RoleType: RoleExecutor},
	}
	svc := NewService(&mockRepo{
		listRolesFn: func(ctx context.Context) ([]Role, error) {
			return expectedRoles, nil
		},
	}, "secret")

	roles, err := svc.ListRoles(context.Background())
	if err != nil {
		t.Fatalf("ListRoles() error = %v", err)
	}
	if len(roles) != 2 {
		t.Errorf("got %d roles, want 2", len(roles))
	}
}
