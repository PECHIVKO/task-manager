package task

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type Repository interface {
	CreateTask(ctx context.Context, task *models.Task) error
	FetchTasks(ctx context.Context, columnID int) ([]*models.Task, error)
	GetTask(ctx context.Context, taskID int) (*models.Task, error)
	DeleteTask(ctx context.Context, taskID int) error
	UpdateTask(ctx context.Context, task *models.Task) error
	ChangeTaskPriority(ctx context.Context, taskID, priority int) error
	MoveToColumn(ctx context.Context, taskID, columnID int) error
}
