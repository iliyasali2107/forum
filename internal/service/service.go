package service

import (
	"forum/internal/repository"
)

type Service struct {
	AuthService
	PostService
	CommentService
	UserService
	ReactionService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		AuthService:     NewAuthService(repository.UserRepository),
		CommentService:  NewCommentService(repository.CommentRepository, repository.UserRepository),
		PostService:     NewPostService(repository.PostRepository, repository.CategoryRepository, repository.UserRepository),
		UserService:     NewUserService(repository.UserRepository),
		ReactionService: NewReactionService(repository.ReactionRepository),
	}
}
