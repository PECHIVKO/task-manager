package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/column"
	"github.com/PECHIVKO/task-manager/models"
)

type ColumnUseCase struct {
	columnRepo column.Repository
}

func NewColumnUseCase(columnRepo column.Repository) *ColumnUseCase {
	return &ColumnUseCase{
		columnRepo: columnRepo,
	}
}

func (c ColumnUseCase) CreateColumn(ctx context.Context, name string, projectID int) error {
	col := &models.Column{
		Name:    name,
		Project: projectID,
	}
	return c.columnRepo.CreateColumn(ctx, col)
}

func (c ColumnUseCase) UpdateColumnName(ctx context.Context, name string, id int) error {

	col := &models.Column{
		ID:   id,
		Name: name,
	}
	return c.columnRepo.UpdateColumnName(ctx, col)
}

func (c ColumnUseCase) MoveColumnToPosition(ctx context.Context, id, position int) error {
	return c.columnRepo.MoveColumnToPosition(ctx, id, position)
}

func (c ColumnUseCase) DeleteColumn(ctx context.Context, id int) error {
	return c.columnRepo.DeleteColumn(ctx, id)
}

func (c ColumnUseCase) GetColumn(ctx context.Context, id string) (*models.Column, error) {
	return c.columnRepo.GetColumn(ctx, id)
}

func (c ColumnUseCase) FetchColumns(ctx context.Context, projectID string) ([]*models.Column, error) {
	return c.columnRepo.FetchColumns(ctx, projectID)
}
