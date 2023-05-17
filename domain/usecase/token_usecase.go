package usecase

import (
	"forum/domain/models"
	"forum/domain/repository"
)

type tokenUsecase struct {
	UserRepository repository.UserRepository
}

type TokenUsecase interface {
	ParseToken(token string) (*models.User, error)
	DeleteToken(token string) error
}

func NewTokenUsecae(UserRepoistory repository.UserRepository) TokenUsecase {
	return &tokenUsecase{
		UserRepository: UserRepoistory,
	}
}

func (tu *tokenUsecase) ParseToken(token string) (*models.User, error) {
	return tu.UserRepository.GetUserByToken(token)
}

func (tu *tokenUsecase) DeleteToken(token string) error {
	return tu.UserRepository.DeleteToken(token)
}
