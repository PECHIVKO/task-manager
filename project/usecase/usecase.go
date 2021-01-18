package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/PECHIVKO/task-manager/project"
)

type ProjectUseCase struct {
	projectRepo project.Repository
}

func NewProjectUseCase(projectRepo project.Repository) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepo: projectRepo,
	}
}

func (p ProjectUseCase) CreateProject(ctx context.Context, name, description string) error {
	pr := &models.Project{
		Name:        name,
		Description: description,
	}
	return p.projectRepo.CreateProject(ctx, pr)
}

func (p ProjectUseCase) UpdateProject(ctx context.Context, name, description string, projectID int) error {

	pr := &models.Project{
		ID:          projectID,
		Name:        name,
		Description: description,
	}
	return p.projectRepo.UpdateProject(ctx, pr)
}

func (p ProjectUseCase) DeleteProject(ctx context.Context, projectID int) error {
	return p.projectRepo.DeleteProject(ctx, projectID)
}

func (p ProjectUseCase) GetProject(ctx context.Context, projectID int) (*models.Project, error) {
	return p.projectRepo.GetProject(ctx, projectID)
}

func (p ProjectUseCase) FetchProjects(ctx context.Context) ([]*models.Project, error) {
	return p.projectRepo.FetchProjects(ctx)
}
