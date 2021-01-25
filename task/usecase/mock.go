package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/stretchr/testify/mock"
)

type TaskUseCaseMock struct {
	mock.Mock
}

func (m TaskUseCaseMock) CreateTask(ctx context.Context, name, description string, columnID int) error {
	args := m.Called(name, description, columnID)

	return args.Error(0)
}

func (m TaskUseCaseMock) UpdateTask(ctx context.Context, name, description string, taskID int) error {
	args := m.Called(name, taskID)

	return args.Error(0)
}

func (m TaskUseCaseMock) MoveToColumn(ctx context.Context, taskID, columnID int) error {
	args := m.Called(taskID, columnID)

	return args.Error(0)
}

func (m TaskUseCaseMock) ChangeTaskPriority(ctx context.Context, taskID, priority int) error {
	args := m.Called(taskID, priority)

	return args.Error(0)
}

func (m TaskUseCaseMock) DeleteTask(ctx context.Context, taskID int) error {
	args := m.Called(taskID)

	return args.Error(0)
}

func (m TaskUseCaseMock) GetTask(ctx context.Context, taskID int) (*models.Task, error) {
	args := m.Called(taskID)

	return args.Get(0).(*models.Task), args.Error(1)
}

func (m TaskUseCaseMock) FetchTasks(ctx context.Context, columnID int) ([]*models.Task, error) {
	args := m.Called(columnID)

	return args.Get(0).([]*models.Task), args.Error(1)
}
