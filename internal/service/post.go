package service

import "forum/internal/repository"

type PostService interface{}

type postService struct {
	Repository repository.PostRepository
}

func NewPostService(repository repository.PostRepository) PostService {
	return &postService{
		Repository: repository,
	}
}
