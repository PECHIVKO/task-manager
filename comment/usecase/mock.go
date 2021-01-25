package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/models"
	"github.com/stretchr/testify/mock"
)

type CommentUseCaseMock struct {
	mock.Mock
}

func (m CommentUseCaseMock) CreateComment(ctx context.Context, text string, taskID int) error {
	args := m.Called(text, taskID)

	return args.Error(0)
}

func (m CommentUseCaseMock) UpdateComment(ctx context.Context, text string, commentID int) error {
	args := m.Called(text, commentID)

	return args.Error(0)
}

func (m CommentUseCaseMock) DeleteComment(ctx context.Context, commentID int) error {
	args := m.Called(commentID)

	return args.Error(0)
}

func (m CommentUseCaseMock) GetComment(ctx context.Context, commentID int) (*models.Comment, error) {
	args := m.Called(commentID)

	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m CommentUseCaseMock) FetchComments(ctx context.Context, taskID int) ([]*models.Comment, error) {
	args := m.Called(taskID)

	return args.Get(0).([]*models.Comment), args.Error(1)
}
