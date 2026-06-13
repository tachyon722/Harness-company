package capability

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type mockCapRepo struct {
	listCapabilitiesFn func(ctx context.Context) ([]Capability, error)
	listBindingsFn     func(ctx context.Context) ([]CapabilityBinding, error)
}

func (m *mockCapRepo) CreateCapability(ctx context.Context, c *Capability) error { return nil }
func (m *mockCapRepo) GetCapability(ctx context.Context, id uuid.UUID) (*Capability, error) {
	return nil, nil
}
func (m *mockCapRepo) ListCapabilities(ctx context.Context) ([]Capability, error) {
	return m.listCapabilitiesFn(ctx)
}
func (m *mockCapRepo) CreateBinding(ctx context.Context, b *CapabilityBinding) error { return nil }
func (m *mockCapRepo) DeleteBinding(ctx context.Context, id uuid.UUID) error { return nil }
func (m *mockCapRepo) ListBindings(ctx context.Context) ([]CapabilityBinding, error) {
	return m.listBindingsFn(ctx)
}
func (m *mockCapRepo) RecordInvocation(ctx context.Context, inv *CapabilityInvocation) error {
	return nil
}

func TestMatchTask_NoCapabilities(t *testing.T) {
	router := NewRouter(&mockCapRepo{
		listCapabilitiesFn: func(ctx context.Context) ([]Capability, error) {
			return []Capability{}, nil
		},
	})

	_, err := router.MatchTask(context.Background(), "do something", "L2")
	if err == nil {
		t.Fatal("expected error when no capabilities exist")
	}
}

func TestMatchTask_ReturnsBestMatch(t *testing.T) {
	router := NewRouter(&mockCapRepo{
		listCapabilitiesFn: func(ctx context.Context) ([]Capability, error) {
			return []Capability{
				{Name: "user_auth", Description: "Authenticate users and manage login sessions", PermissionLevel: "L2"},
				{Name: "data_export", Description: "Export data to CSV and Excel formats", PermissionLevel: "L3"},
				{Name: "email_service", Description: "Send transactional emails", PermissionLevel: "L1"},
			}, nil
		},
	})

	result, err := router.MatchTask(context.Background(), "Need to authenticate a user with email and password", "L2")
	if err != nil {
		t.Fatalf("MatchTask() error = %v", err)
	}

	if result.Capability.Name != "user_auth" {
		t.Errorf("expected user_auth, got %s", result.Capability.Name)
	}
}

func TestMatchTask_FiltersByPermission(t *testing.T) {
	router := NewRouter(&mockCapRepo{
		listCapabilitiesFn: func(ctx context.Context) ([]Capability, error) {
			return []Capability{
				{Name: "admin_panel", Description: "Full admin access to system", PermissionLevel: "L4"},
				{Name: "basic_info", Description: "View basic information", PermissionLevel: "L1"},
			}, nil
		},
	})

	result, err := router.MatchTask(context.Background(), "access the admin panel", "L2")
	if err != nil {
		t.Fatalf("MatchTask() error = %v", err)
	}

	if result.Capability.Name != "basic_info" {
		t.Errorf("expected basic_info for L2 user (L4 too high), got %s", result.Capability.Name)
	}
}

func TestMatchTask_ExactMatchPreferred(t *testing.T) {
	router := NewRouter(&mockCapRepo{
		listCapabilitiesFn: func(ctx context.Context) ([]Capability, error) {
			return []Capability{
				{Name: "search_users", Description: "Search users by name and email", PermissionLevel: "L1"},
				{Name: "user_management", Description: "Full user management including create, update, delete", PermissionLevel: "L2"},
				{Name: "reporting", Description: "Generate reports and analytics", PermissionLevel: "L1"},
			}, nil
		},
	})

	result, err := router.MatchTask(context.Background(), "create a new user account and assign roles", "L2")
	if err != nil {
		t.Fatalf("MatchTask() error = %v", err)
	}

	if result.Capability.Name != "user_management" {
		t.Errorf("expected user_management, got %s", result.Capability.Name)
	}
}

func TestMatchTask_CaseInsensitive(t *testing.T) {
	router := NewRouter(&mockCapRepo{
		listCapabilitiesFn: func(ctx context.Context) ([]Capability, error) {
			return []Capability{
				{Name: "EMAIL", Description: "SEND EMAIL NOTIFICATIONS", PermissionLevel: "L1"},
				{Name: "sms", Description: "send sms messages", PermissionLevel: "L1"},
			}, nil
		},
	})

	result, err := router.MatchTask(context.Background(), "send email to user", "L1")
	if err != nil {
		t.Fatalf("MatchTask() error = %v", err)
	}

	if result.Capability.Name != "EMAIL" {
		t.Errorf("expected EMAIL (case-insensitive match), got %s", result.Capability.Name)
	}
}
