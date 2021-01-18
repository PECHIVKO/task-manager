package task

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type UseCase interface {
	CreateTask(ctx context.Context, name, description string, columnID int) error
	FetchTasks(ctx context.Context, columnID int) ([]*models.Task, error)
	GetTask(ctx context.Context, taskID int) (*models.Task, error)
	DeleteTask(ctx context.Context, taskID int) error
	UpdateTask(ctx context.Context, name, description string, taskID int) error
	ChangeTaskPriority(ctx context.Context, taskID, priority int) error
	MoveToColumn(ctx context.Context, taskID, columnID int) error
}
