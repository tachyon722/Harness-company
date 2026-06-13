package verification

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

type mockVerRepo struct {
	createReportFn func(ctx context.Context, report *VerificationReport) error
	getReportFn    func(ctx context.Context, id uuid.UUID) (*VerificationReport, error)
	listReportsFn  func(ctx context.Context, workflowID string) ([]VerificationReport, error)
	createReviewFn func(ctx context.Context, review *ReviewAssignment) error
	updateReviewFn func(ctx context.Context, id uuid.UUID, result map[string]any, status string) error
}

func (m *mockVerRepo) CreateReport(ctx context.Context, report *VerificationReport) error {
	return m.createReportFn(ctx, report)
}
func (m *mockVerRepo) GetReport(ctx context.Context, id uuid.UUID) (*VerificationReport, error) {
	return m.getReportFn(ctx, id)
}
func (m *mockVerRepo) ListReports(ctx context.Context, workflowID string) ([]VerificationReport, error) {
	return m.listReportsFn(ctx, workflowID)
}
func (m *mockVerRepo) CreateReview(ctx context.Context, review *ReviewAssignment) error {
	return m.createReviewFn(ctx, review)
}
func (m *mockVerRepo) UpdateReview(ctx context.Context, id uuid.UUID, result map[string]any, status string) error {
	return m.updateReviewFn(ctx, id, result, status)
}

func TestCreateReport_ComputesScore(t *testing.T) {
	var saved *VerificationReport
	svc := NewService(&mockVerRepo{
		createReportFn: func(ctx context.Context, report *VerificationReport) error {
			saved = report
			return nil
		},
	})

	report, err := svc.CreateReport(context.Background(), "wf-1", "task-1", 85.0, 90.0, 80.0, nil, nil)
	if err != nil {
		t.Fatalf("CreateReport() error = %v", err)
	}

	expectedOverall := 85*0.4 + 90*0.35 + 80*0.25
	if report.OverallScore != expectedOverall {
		t.Errorf("OverallScore = %f, want %f", report.OverallScore, expectedOverall)
	}
}

func TestCreateReport_RoundToTwoDecimals(t *testing.T) {
	var saved *VerificationReport
	svc := NewService(&mockVerRepo{
		createReportFn: func(ctx context.Context, report *VerificationReport) error {
			saved = report
			return nil
		},
	})

	report, err := svc.CreateReport(context.Background(), "wf-1", "task-1", 33.33, 66.67, 50.0, nil, nil)
	if err != nil {
		t.Fatalf("CreateReport() error = %v", err)
	}

	expectedOverall := 33.33*0.4 + 66.67*0.35 + 50.0*0.25
	if report.OverallScore != expectedOverall {
		t.Errorf("OverallScore = %f, want %f", report.OverallScore, expectedOverall)
	}
}

func TestAssignReview_InvalidLevel(t *testing.T) {
	svc := NewService(&mockVerRepo{})

	_, err := svc.AssignReview(context.Background(), uuid.New(), "L4", "machine")
	if err == nil {
		t.Fatal("expected error for invalid level L4")
	}
}

func TestAssignReview_InvalidReviewer(t *testing.T) {
	svc := NewService(&mockVerRepo{})

	_, err := svc.AssignReview(context.Background(), uuid.New(), "L1", "unknown")
	if err == nil {
		t.Fatal("expected error for unknown reviewer type")
	}
}

func TestAssignReview_Success(t *testing.T) {
	var saved *ReviewAssignment
	svc := NewService(&mockVerRepo{
		createReviewFn: func(ctx context.Context, review *ReviewAssignment) error {
			saved = review
			return nil
		},
	})

	reportID := uuid.New()
	review, err := svc.AssignReview(context.Background(), reportID, "L2", "ai")
	if err != nil {
		t.Fatalf("AssignReview() error = %v", err)
	}

	if review.Level != "L2" {
		t.Errorf("Level = %s, want L2", review.Level)
	}
	if review.ReviewerType != "ai" {
		t.Errorf("ReviewerType = %s, want ai", review.ReviewerType)
	}
	if review.ReportID != reportID {
		t.Errorf("ReportID = %s, want %s", review.ReportID, reportID)
	}
	if review.Status != "pending" {
		t.Errorf("Status = %s, want pending", review.Status)
	}
}
