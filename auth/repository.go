package auth

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, username string) (*models.User, error)
}
