package service

import (
	"forum/internal/repository"
)

type Service struct {
	AuthService
	PostService
	CommentService
	UserService
	VoteService
}

func NewService(repository *repository.Repository) *Service {
	return &Service{
		AuthService:    NewAuthService(repository.UserRepository),
		CommentService: NewCommentService(repository.CommentRepository),
		PostService:    NewPostService(repository.PostRepository, repository.CategoryRepository),
		UserService:    NewUserService(repository.UserRepository),
		VoteService:    NewVoteService(repository.VoteRepository),
	}
}
