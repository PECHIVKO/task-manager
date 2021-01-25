package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/stretchr/testify/mock"
)

type ProjectUseCaseMock struct {
	mock.Mock
}

func (m ProjectUseCaseMock) CreateProject(ctx context.Context, name, description string) error {
	args := m.Called(name, description)

	return args.Error(0)
}

func (m ProjectUseCaseMock) UpdateProject(ctx context.Context, name, description string, projectID int) error {
	args := m.Called(name, description, projectID)

	return args.Error(0)
}

func (m ProjectUseCaseMock) DeleteProject(ctx context.Context, columnID int) error {
	args := m.Called(columnID)

	return args.Error(0)
}

func (m ProjectUseCaseMock) GetProject(ctx context.Context, projectID int) (*models.Project, error) {
	args := m.Called(projectID)

	return args.Get(0).(*models.Project), args.Error(1)
}

func (m ProjectUseCaseMock) FetchProjects(ctx context.Context) ([]*models.Project, error) {
	args := m.Called()

	return args.Get(0).([]*models.Project), args.Error(1)
}
