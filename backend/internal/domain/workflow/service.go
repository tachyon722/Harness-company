package workflow

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

func (s *Service) CreateTemplate(ctx context.Context, input CreateWorkflowInput) (*WorkflowTemplate, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}
	if len(input.Stages) == 0 {
		input.Stages = defaultStages()
	}
	if input.AssigneeType == "" {
		input.AssigneeType = "either"
	}
	return s.repo.CreateTemplate(ctx, input)
}

func (s *Service) GetTemplate(ctx context.Context, id uuid.UUID) (*WorkflowTemplate, error) {
	return s.repo.GetTemplate(ctx, id)
}

func (s *Service) ListTemplates(ctx context.Context) ([]WorkflowTemplate, error) {
	return s.repo.ListTemplates(ctx)
}

func (s *Service) StartWorkflow(ctx context.Context, input StartWorkflowInput) (*WorkflowInstance, error) {
	tmpl, err := s.repo.GetTemplate(ctx, input.TemplateID)
	if err != nil {
		return nil, fmt.Errorf("template not found: %w", err)
	}
	if input.Context == nil {
		input.Context = map[string]any{}
	}

	inst, err := s.repo.CreateInstance(ctx, input)
	if err != nil {
		return nil, err
	}

	for i, stage := range tmpl.Stages {
		task := &Task{
			WorkflowID:     inst.ID,
			Stage:          i,
			StageType:      stage.Type,
			AssigneeType:   stage.AssigneeType,
			Input:          input.Context,
			WeightSnapshot: tmpl.RequiredWeight,
			Status:         TaskPending,
		}
		if i == 0 {
			task.Status = TaskAssigned
		}
		if _, err := s.repo.CreateTask(ctx, task); err != nil {
			return nil, fmt.Errorf("create task for stage %d: %w", i, err)
		}
	}

	return inst, nil
}

func (s *Service) GetWorkflow(ctx context.Context, id uuid.UUID) (*WorkflowInstance, error) {
	inst, err := s.repo.GetInstance(ctx, id)
	if err != nil {
		return nil, err
	}
	tasks, err := s.repo.GetTasksByWorkflow(ctx, id)
	if err != nil {
		return nil, err
	}
	inst.Tasks = tasks
	return inst, nil
}

func (s *Service) CompleteTask(ctx context.Context, taskID uuid.UUID, output map[string]any) error {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}
	if err := s.repo.UpdateTaskStatus(ctx, taskID, TaskCompleted, output); err != nil {
		return err
	}

	inst, err := s.repo.GetInstance(ctx, task.WorkflowID)
	if err != nil {
		return err
	}

	nextStage := inst.CurrentStage + 1
	if err := s.repo.UpdateInstanceStage(ctx, inst.ID, nextStage); err != nil {
		return err
	}

	tmpl, err := s.repo.GetTemplate(ctx, inst.TemplateID)
	if err != nil {
		return err
	}
	if nextStage >= len(tmpl.Stages) {
		return s.repo.UpdateInstanceStatus(ctx, inst.ID, WorkflowCompleted)
	}

	return nil
}

func (s *Service) RecordDecision(ctx context.Context, taskID uuid.UUID, decisionMakerID uuid.UUID, makerType string, reasoning string, outcome string, input, output map[string]any) (*Decision, error) {
	d := &Decision{
		TaskID:          taskID,
		DecisionMakerID: decisionMakerID,
		MakerType:       makerType,
		Weight:          1.0,
		Input:           input,
		Output:          output,
		Reasoning:       reasoning,
		Outcome:         outcome,
	}
	return s.repo.RecordDecision(ctx, d)
}

func (s *Service) GetContext(ctx context.Context, workflowID uuid.UUID) (*WorkflowContext, error) {
	return s.repo.GetWorkflowContext(ctx, workflowID)
}

func (s *Service) UpdateContext(ctx context.Context, wc *WorkflowContext) error {
	return s.repo.UpsertWorkflowContext(ctx, wc)
}

func defaultStages() []Stage {
	return []Stage{
		{Type: StagePlan, Name: "Planning", AssigneeType: "either"},
		{Type: StageExecute, Name: "Execution", AssigneeType: "either"},
		{Type: StageReview, Name: "Review", AssigneeType: "either"},
	}
}
