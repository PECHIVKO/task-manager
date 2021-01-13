package column

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type UseCase interface {
	CreateColumn(ctx context.Context, name string, project_id int) error
	FetchColumns(ctx context.Context, projectID string) ([]*models.Column, error)
	GetColumn(ctx context.Context, columnID string) (*models.Column, error)
	DeleteColumn(ctx context.Context, columnID int) error
	UpdateColumnName(ctx context.Context, name string, id int) error
	MoveColumnToPosition(ctx context.Context, id, position int) error
}
