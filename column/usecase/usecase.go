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

func (c ColumnUseCase) UpdateColumnName(ctx context.Context, name string, columnID int) error {

	col := &models.Column{
		ID:   columnID,
		Name: name,
	}
	return c.columnRepo.UpdateColumnName(ctx, col)
}

func (c ColumnUseCase) MoveColumnToPosition(ctx context.Context, columnID, position int) error {
	return c.columnRepo.MoveColumnToPosition(ctx, columnID, position)
}

func (c ColumnUseCase) DeleteColumn(ctx context.Context, columnID int) error {
	return c.columnRepo.DeleteColumn(ctx, columnID)
}

func (c ColumnUseCase) GetColumn(ctx context.Context, columnID int) (*models.Column, error) {
	return c.columnRepo.GetColumn(ctx, columnID)
}

func (c ColumnUseCase) FetchColumns(ctx context.Context, projectID int) ([]*models.Column, error) {
	return c.columnRepo.FetchColumns(ctx, projectID)
}
