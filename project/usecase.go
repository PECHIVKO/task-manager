package project

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type UseCase interface {
	CreateProject(ctx context.Context, name, description string) error
	FetchProjects(ctx context.Context) ([]*models.Project, error)
	GetProject(ctx context.Context, projectID int) (*models.Project, error)
	DeleteProject(ctx context.Context, projectID int) error
	UpdateProject(ctx context.Context, name, description string, projectID int) error
}
