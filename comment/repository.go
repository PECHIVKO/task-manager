package comment

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type Repository interface {
	CreateComment(ctx context.Context, comment *models.Comment) error
	FetchComments(ctx context.Context, taskID int) ([]*models.Comment, error)
	GetComment(ctx context.Context, commentID int) (*models.Comment, error)
	DeleteComment(ctx context.Context, commentID int) error
	UpdateComment(ctx context.Context, comment *models.Comment) error
}
