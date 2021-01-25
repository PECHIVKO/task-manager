package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/stretchr/testify/mock"
)

type ColumnUseCaseMock struct {
	mock.Mock
}

func (m ColumnUseCaseMock) CreateColumn(ctx context.Context, name string, projectID int) error {
	args := m.Called(name, projectID)

	return args.Error(0)
}

func (m ColumnUseCaseMock) UpdateColumnName(ctx context.Context, name string, columnID int) error {
	args := m.Called(name, columnID)

	return args.Error(0)
}

func (m ColumnUseCaseMock) MoveColumnToPosition(ctx context.Context, columnID, position int) error {
	args := m.Called(columnID, position)

	return args.Error(0)
}

func (m ColumnUseCaseMock) DeleteColumn(ctx context.Context, columnID int) error {
	args := m.Called(columnID)

	return args.Error(0)
}

func (m ColumnUseCaseMock) GetColumn(ctx context.Context, columnID int) (*models.Column, error) {
	args := m.Called(columnID)

	return args.Get(0).(*models.Column), args.Error(1)
}

func (m ColumnUseCaseMock) FetchColumns(ctx context.Context, projectID int) ([]*models.Column, error) {
	args := m.Called(projectID)

	return args.Get(0).([]*models.Column), args.Error(1)
}
