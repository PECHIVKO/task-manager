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

func (p ProjectUseCase) UpdateProject(ctx context.Context, name, description string, id int) error {

	pr := &models.Project{
		ID:          id,
		Name:        name,
		Description: description,
	}
	return p.projectRepo.UpdateProject(ctx, pr)
}

func (p ProjectUseCase) DeleteProject(ctx context.Context, id string) error {
	return p.projectRepo.DeleteProject(ctx, id)
}

func (p ProjectUseCase) GetProject(ctx context.Context, id string) (*models.Project, error) {
	return p.projectRepo.GetProject(ctx, id)
}

func (p ProjectUseCase) FetchProjects(ctx context.Context) ([]*models.Project, error) {
	return p.projectRepo.FetchProjects(ctx)
}
