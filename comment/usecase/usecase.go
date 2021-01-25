package usecase

import (
	"context"

	"github.com/PECHIVKO/task-manager/comment"
	"github.com/PECHIVKO/task-manager/models"
)

type CommentUseCase struct {
	commentRepo comment.Repository
}

func NewCommentUseCase(commentRepo comment.Repository) *CommentUseCase {
	return &CommentUseCase{
		commentRepo: commentRepo,
	}
}

func (c CommentUseCase) CreateComment(ctx context.Context, text string, taskID int) error {
	tsk := &models.Comment{
		Comment: text,
		Task:    taskID,
	}
	return c.commentRepo.CreateComment(ctx, tsk)
}

func (c CommentUseCase) UpdateComment(ctx context.Context, text string, commentID int) error {

	tsk := &models.Comment{
		ID:      commentID,
		Comment: text,
	}
	return c.commentRepo.UpdateComment(ctx, tsk)
}

func (c CommentUseCase) DeleteComment(ctx context.Context, commentID int) error {
	return c.commentRepo.DeleteComment(ctx, commentID)
}

func (c CommentUseCase) GetComment(ctx context.Context, commentID int) (*models.Comment, error) {
	return c.commentRepo.GetComment(ctx, commentID)
}

func (c CommentUseCase) FetchComments(ctx context.Context, taskID int) ([]*models.Comment, error) {
	return c.commentRepo.FetchComments(ctx, taskID)
}
