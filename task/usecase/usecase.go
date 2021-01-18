package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/PECHIVKO/task-manager/task"
)

type TaskUseCase struct {
	taskRepo task.Repository
}

func NewTaskUseCase(taskRepo task.Repository) *TaskUseCase {
	return &TaskUseCase{
		taskRepo: taskRepo,
	}
}

func (c TaskUseCase) CreateTask(ctx context.Context, name, description string, columnID int) error {
	tsk := &models.Task{
		Name:        name,
		Description: description,
		Column:      columnID,
	}
	return c.taskRepo.CreateTask(ctx, tsk)
}

func (c TaskUseCase) UpdateTask(ctx context.Context, name, description string, taskID int) error {

	tsk := &models.Task{
		ID:          taskID,
		Name:        name,
		Description: description,
	}
	return c.taskRepo.UpdateTask(ctx, tsk)
}

func (c TaskUseCase) MoveToColumn(ctx context.Context, taskID, columnID int) error {
	return c.taskRepo.MoveToColumn(ctx, taskID, columnID)
}

func (c TaskUseCase) DeleteTask(ctx context.Context, taskID int) error {
	return c.taskRepo.DeleteTask(ctx, taskID)
}

func (c TaskUseCase) GetTask(ctx context.Context, taskID int) (*models.Task, error) {
	return c.taskRepo.GetTask(ctx, taskID)
}

func (c TaskUseCase) FetchTasks(ctx context.Context, columnID int) ([]*models.Task, error) {
	return c.taskRepo.FetchTasks(ctx, columnID)
}

func (c TaskUseCase) ChangeTaskPriority(ctx context.Context, taskID, priority int) error {
	return c.taskRepo.ChangeTaskPriority(ctx, taskID, priority)
}
