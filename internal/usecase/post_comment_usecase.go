package usecase

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg/validator"
)

type postCommentUsecase struct {
	commentRepository repository.CommentRepository
}

type PostCommentUsecase interface {
	CreateComment(*validator.Validator, *models.Comment) error
}

func NewPostCommentUsecase(commentRepository repository.CommentRepository, userRepository repository.UserRepository) PostCommentUsecase {
	return &postCommentUsecase{
		commentRepository: commentRepository,
	}
}

func (pcu *postCommentUsecase) CreateComment(v *validator.Validator, comment *models.Comment) error {
	if validateComment(v, comment); !v.Valid() {
		return ErrFormValidation
	}

	_, err := pcu.commentRepository.CreateComment(comment)
	if err != nil {
		return ErrInternalServer
	}

	return nil
}
