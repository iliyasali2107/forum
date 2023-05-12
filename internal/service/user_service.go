package service

import "forum/internal/repository"

type UserService interface{}

type userService struct {
	Repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) UserService {
	return &userService{
		Repository: repository,
	}
}
