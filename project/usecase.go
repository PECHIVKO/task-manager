package project

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type UseCase interface {
	CreateProject(ctx context.Context, name, description string) error
	FetchProjects(ctx context.Context) ([]*models.Project, error)
	GetProject(ctx context.Context, projectID string) (*models.Project, error)
	DeleteProject(ctx context.Context, projectID string) error
	UpdateProject(ctx context.Context, name, description string, id int) error
}
