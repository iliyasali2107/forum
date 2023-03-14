package service

import "forum/internal/repository"

type CommentService interface{}

type commentService struct {
	Repository repository.CommentRepository
}

func NewCommentService(repository repository.CommentRepository) CommentService {
	return &commentService{
		Repository: repository,
	}
}
