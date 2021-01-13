package usecase

import (
	"context"
	"fmt"

	"github.com/PECHIVKO/task-manager/auth"
	"github.com/PECHIVKO/task-manager/models"
)

type AuthUseCase struct {
	userRepo auth.UserRepository
}

func NewAuthUseCase(userRepo auth.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, username string) error {

	user := &models.User{
		Username: username,
	}

	return a.userRepo.CreateUser(ctx, user)
}

func (a *AuthUseCase) SignIn(ctx context.Context, username string) error {

	user, err := a.userRepo.GetUser(ctx, username)
	if err != nil {
		return fmt.Errorf("signIn error")
		//return auth.ErrUserNotFound
	}

	ctx = context.WithValue(context.Background(), auth.CtxUserKey, user.Username)

	return nil
}
