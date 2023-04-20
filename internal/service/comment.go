package service

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg/validator"
)

type CommentService interface {
	CreateComment(*validator.Validator, *models.Comment) error
}

type commentService struct {
	cr repository.CommentRepository
	ur repository.UserRepository
}

func NewCommentService(commentRepository repository.CommentRepository, userRepository repository.UserRepository) CommentService {
	return &commentService{
		cr: commentRepository,
		ur: userRepository,
	}
}

func (cs *commentService) CreateComment(v *validator.Validator, comment *models.Comment) error {
	if validateComment(v, comment); !v.Valid() {
		return ErrFormValidation
	}

	_, err := cs.cr.CreateCommentWithoutParent(comment)
	if err != nil {
		return err
	}

	return nil
}

func validateComment(v *validator.Validator, comment *models.Comment) {
	v.Check(comment.Content != "", "content", "must be provided")
}
