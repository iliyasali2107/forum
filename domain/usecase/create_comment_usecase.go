package usecase

import (
	"fmt"
	"time"

	"forum/domain/models"
	"forum/domain/repository"
)

type createCommentUsecase struct {
	commentRepository repository.CommentRepository
	ContextTimeout    time.Duration
}

type CreateCommentUsecase interface {
	CreateComment(*models.Comment) error
}

func NewCreateCommentUsecase(commentRepository repository.CommentRepository, userRepository repository.UserRepository, timeout time.Duration) CreateCommentUsecase {
	return &createCommentUsecase{
		commentRepository: commentRepository,
		ContextTimeout:    timeout,
	}
}

func (pcu *createCommentUsecase) CreateComment(comment *models.Comment) error {

	_, err := pcu.commentRepository.CreateComment(comment)
	if err != nil {
		return fmt.Errorf("couldn't create comment: %w", err)
	}

	return nil
}
