package layer

import (
	"context"
	"testing"
)

type mockLayerRepo struct {
	getConfigFn func(ctx context.Context, mvruID string, layer LayerType) (*LayerConfig, error)
}

func (m *mockLayerRepo) GetConfig(ctx context.Context, mvruID string, layer LayerType) (*LayerConfig, error) {
	return m.getConfigFn(ctx, mvruID, layer)
}
func (m *mockLayerRepo) UpsertConfig(ctx context.Context, cfg *LayerConfig) error { return nil }
func (m *mockLayerRepo) ListRoutingRules(ctx context.Context) ([]LayerRoutingRule, error) {
	return nil, nil
}

func TestClassify_SimpleTask(t *testing.T) {
	svc := NewClassifierService(&mockLayerRepo{})

	result, err := svc.Classify(context.Background(), "fix login button css", "operational")
	if err != nil {
		t.Fatalf("Classify() error = %v", err)
	}
	if result.Layer != LayerOperational {
		t.Errorf("expected operational, got %s", result.Layer)
	}
	if result.Confidence <= 0 || result.Confidence > 1 {
		t.Errorf("confidence out of range: %f", result.Confidence)
	}
}

func TestClassify_StrategicTask(t *testing.T) {
	svc := NewClassifierService(&mockLayerRepo{})

	desc := `Long range strategic plan for organizational transformation including market expansion, merger acquisition targets, and five year revenue growth forecasting with multiple scenario analyses and risk assessment frameworks`
	result, err := svc.Classify(context.Background(), desc, "")
	if err != nil {
		t.Fatalf("Classify() error = %v", err)
	}
	if result.Layer != LayerStrategic {
		t.Errorf("expected strategic, got %s", result.Layer)
	}
}

func TestClassify_OverridesLayer(t *testing.T) {
	svc := NewClassifierService(&mockLayerRepo{})

	result, err := svc.Classify(context.Background(), "any task", "strategic")
	if err != nil {
		t.Fatalf("Classify() error = %v", err)
	}
	if result.Layer != LayerStrategic {
		t.Errorf("expected strategic, got %s", result.Layer)
	}
}

func TestClassify_RiskKeywords(t *testing.T) {
	svc := NewClassifierService(&mockLayerRepo{})

	result, err := svc.Classify(context.Background(), "Handle urgent security vulnerability compliance audit with critical risk assessment", "")
	if err != nil {
		t.Fatalf("Classify() error = %v", err)
	}
	if result.Layer != LayerTactical {
		t.Errorf("expected tactical for urgent/risk task, got %s", result.Layer)
	}
}
