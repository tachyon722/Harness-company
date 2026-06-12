package governance

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrValidation = errors.New("validation error")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePermission(ctx context.Context, p *Permission) (*Permission, error) {
	if p.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if p.Level < 1 || p.Level > 4 {
		return nil, fmt.Errorf("%w: level must be between 1 and 4", ErrValidation)
	}
	if p.Behavior != "auto" && p.Behavior != "notify" && p.Behavior != "approve" && p.Behavior != "deny" {
		p.Behavior = "notify"
	}
	return s.repo.CreatePermission(ctx, p)
}

func (s *Service) ListPermissions(ctx context.Context) ([]Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *Service) CreatePrinciple(ctx context.Context, input CreatePrincipleInput) (*Principle, error) {
	if input.Name == "" || input.Description == "" {
		return nil, fmt.Errorf("%w: name and description are required", ErrValidation)
	}
	return s.repo.CreatePrinciple(ctx, input)
}

func (s *Service) ListPrinciples(ctx context.Context) ([]Principle, error) {
	return s.repo.ListPrinciples(ctx)
}

func (s *Service) GetPrinciple(ctx context.Context, id uuid.UUID) (*Principle, error) {
	return s.repo.GetPrinciple(ctx, id)
}

func (s *Service) CreateControlRule(ctx context.Context, input CreateControlRuleInput) (*ControlRule, error) {
	if input.Action == "" {
		return nil, fmt.Errorf("%w: action is required", ErrValidation)
	}
	return s.repo.CreateControlRule(ctx, input)
}

func (s *Service) ListControlRules(ctx context.Context) ([]ControlRule, error) {
	return s.repo.ListControlRules(ctx)
}

func (s *Service) CheckPermission(ctx context.Context, input PermissionCheckInput) (*PermissionCheckResult, error) {
	permissions, err := s.repo.ListPermissions(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range permissions {
		switch p.Behavior {
		case "auto":
			return &PermissionCheckResult{Allowed: true, Level: p.Level, Behavior: p.Behavior, Reason: "auto-allowed"}, nil
		case "notify":
			return &PermissionCheckResult{Allowed: true, Level: p.Level, Behavior: p.Behavior, Reason: "notify-allowed"}, nil
		case "approve":
			return &PermissionCheckResult{Allowed: false, Level: p.Level, Behavior: p.Behavior, Reason: "requires approval"}, nil
		case "deny":
			return &PermissionCheckResult{Allowed: false, Level: p.Level, Behavior: p.Behavior, Reason: "denied"}, nil
		}
	}

	return &PermissionCheckResult{Allowed: false, Level: 0, Behavior: "", Reason: "no matching permission"}, nil
}
