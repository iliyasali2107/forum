package usecase

import (
	"time"

	"forum/domain/models"
	"forum/domain/repository"
	"forum/pkg/validator"
)

type createCommentUsecase struct {
	commentRepository repository.CommentRepository
	ContextTimeout    time.Duration
}

type CreateCommentUsecase interface {
	CreateComment(*validator.Validator, *models.Comment) error
}

func NewCreateCommentUsecase(commentRepository repository.CommentRepository, userRepository repository.UserRepository, timeout time.Duration) CreateCommentUsecase {
	return &createCommentUsecase{
		commentRepository: commentRepository,
		ContextTimeout:    timeout,
	}
}

func (pcu *createCommentUsecase) CreateComment(v *validator.Validator, comment *models.Comment) error {
	if validateComment(v, comment); !v.Valid() {
		return ErrFormValidation
	}

	_, err := pcu.commentRepository.CreateComment(comment)
	if err != nil {
		return ErrInternalServer
	}

	return nil
}
