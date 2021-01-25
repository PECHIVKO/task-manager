package column

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type UseCase interface {
	CreateColumn(ctx context.Context, name string, projectID int) error
	FetchColumns(ctx context.Context, projectID int) ([]*models.Column, error)
	GetColumn(ctx context.Context, columnID int) (*models.Column, error)
	DeleteColumn(ctx context.Context, columnID int) error
	UpdateColumnName(ctx context.Context, name string, columnID int) error
	MoveColumnToPosition(ctx context.Context, id, position int) error
}
