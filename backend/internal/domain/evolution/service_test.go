package evolution

import (
	"context"
	"testing"
)

type mockEvoRepo struct {
	getWeightsFn   func(ctx context.Context, actorType, actorID string) (*DecisionWeight, error)
	upsertWeightFn func(ctx context.Context, w *DecisionWeight) error
	getAlphasFn    func(ctx context.Context) (*AlphaConfig, error)
}

func (m *mockEvoRepo) GetWeight(ctx context.Context, actorType, actorID string) (*DecisionWeight, error) {
	return m.getWeightsFn(ctx, actorType, actorID)
}
func (m *mockEvoRepo) UpsertWeight(ctx context.Context, w *DecisionWeight) error {
	return m.upsertWeightFn(ctx, w)
}
func (m *mockEvoRepo) GetAlphas(ctx context.Context) (*AlphaConfig, error) {
	return m.getAlphasFn(ctx)
}
func (m *mockEvoRepo) UpsertAlphas(ctx context.Context, a *AlphaConfig) error { return nil }
func (m *mockEvoRepo) ListWeights(ctx context.Context, actorType string) ([]DecisionWeight, error) {
	return nil, nil
}
func (m *mockEvoRepo) ListExperiments(ctx context.Context, mvruID string) ([]Experiment, error) {
	return nil, nil
}
func (m *mockEvoRepo) CreateExperiment(ctx context.Context, e *Experiment) error { return nil }
func (m *mockEvoRepo) UpdateExperiment(ctx context.Context, e *Experiment) error { return nil }
func (m *mockEvoRepo) CreateKnowledge(ctx context.Context, k *KnowledgeEntry) error { return nil }
func (m *mockEvoRepo) ListKnowledge(ctx context.Context, tags []string) ([]KnowledgeEntry, error) {
	return nil, nil
}
func (m *mockEvoRepo) CreateSignal(ctx context.Context, s *Signal) error { return nil }
func (m *mockEvoRepo) ListSignals(ctx context.Context, signalType string, acknowledged *bool) ([]Signal, error) {
	return nil, nil
}
func (m *mockEvoRepo) AcknowledgeSignal(ctx context.Context, id string) error { return nil }

func TestComputeWeight_WithDefaultAlphas(t *testing.T) {
	svc := NewService(&mockEvoRepo{
		getAlphasFn: func(ctx context.Context) (*AlphaConfig, error) {
			return &AlphaConfig{
				ExpertiseAlpha:   0.2,
				TrackRecordAlpha: 0.2,
				ReliabilityAlpha: 0.2,
				RecencyAlpha:     0.1,
				ContextFitAlpha:  0.15,
				PrincipleAlpha:   0.15,
			}, nil
		},
	})

	weight := &DecisionWeight{
		ActorType:       "human",
		ActorID:         "user-1",
		Expertise:       80,
		TrackRecord:     70,
		Reliability:     90,
		Recency:         60,
		ContextFit:      75,
		PrincipleScore:  85,
	}

	result, err := svc.ComputeWeight(context.Background(), weight)
	if err != nil {
		t.Fatalf("ComputeWeight() error = %v", err)
	}

	expected := 80*0.2 + 70*0.2 + 90*0.2 + 60*0.1 + 75*0.15 + 85*0.15
	if result.OverallScore != expected {
		t.Errorf("OverallScore = %f, want %f", result.OverallScore, expected)
	}
}

func TestComputeWeight_PerfectScores(t *testing.T) {
	svc := NewService(&mockEvoRepo{
		getAlphasFn: func(ctx context.Context) (*AlphaConfig, error) {
			return &AlphaConfig{
				ExpertiseAlpha:   0.2,
				TrackRecordAlpha: 0.2,
				ReliabilityAlpha: 0.2,
				RecencyAlpha:     0.1,
				ContextFitAlpha:  0.15,
				PrincipleAlpha:   0.15,
			}, nil
		},
	})

	weight := &DecisionWeight{
		ActorType:      "ai",
		ActorID:        "agent-1",
		Expertise:      100,
		TrackRecord:    100,
		Reliability:    100,
		Recency:        100,
		ContextFit:     100,
		PrincipleScore: 100,
	}

	result, err := svc.ComputeWeight(context.Background(), weight)
	if err != nil {
		t.Fatalf("ComputeWeight() error = %v", err)
	}

	if result.OverallScore != 100 {
		t.Errorf("OverallScore = %f, want 100", result.OverallScore)
	}
}

func TestRecordOutcome(t *testing.T) {
	var saved *DecisionWeight
	svc := NewService(&mockEvoRepo{
		getWeightsFn: func(ctx context.Context, actorType, actorID string) (*DecisionWeight, error) {
			return &DecisionWeight{
				ActorType:      actorType,
				ActorID:        actorID,
				Expertise:      50,
				TrackRecord:    50,
				Reliability:    50,
				Recency:        50,
				ContextFit:     50,
				PrincipleScore: 50,
				OverallScore:   50,
			}, nil
		},
		upsertWeightFn: func(ctx context.Context, w *DecisionWeight) error {
			saved = w
			return nil
		},
	})

	err := svc.RecordOutcome(context.Background(), "human", "user-1", "success", 90)
	if err != nil {
		t.Fatalf("RecordOutcome() error = %v", err)
	}

	if saved == nil {
		t.Fatal("expected weight to be saved")
	}
	if saved.TrackRecord <= 50 {
		t.Errorf("expected TrackRecord to increase, got %f", saved.TrackRecord)
	}
}

func TestValidateAlphas_SumToApproxOne(t *testing.T) {
	tests := []struct {
		name   string
		config AlphaConfig
		valid  bool
	}{
		{"sum = 1.0", AlphaConfig{0.2, 0.2, 0.2, 0.1, 0.15, 0.15}, true},
		{"sum = 1.0 exact", AlphaConfig{0.1, 0.1, 0.1, 0.2, 0.2, 0.3}, true},
		{"sum too low", AlphaConfig{0.1, 0.1, 0.1, 0.1, 0.1, 0.1}, false},
		{"sum too high", AlphaConfig{0.3, 0.3, 0.3, 0.3, 0.3, 0.3}, false},
		{"sum approx 1.0", AlphaConfig{0.3, 0.3, 0.2, 0.1, 0.05, 0.05}, true},
	}

	svc := NewService(&mockEvoRepo{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.UpdateAlphas(context.Background(), &tt.config)
			if tt.valid && err != nil {
				t.Errorf("expected valid, got error: %v", err)
			}
			if !tt.valid && err == nil {
				t.Errorf("expected invalid, got no error")
			}
		})
	}
}
