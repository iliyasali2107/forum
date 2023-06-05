package usecase

import (
	"time"

	"forum/domain/models"
	"forum/domain/repository"
	"forum/pkg/utils"
)

type signupUsecase struct {
	userRepository repository.UserRepository
	ContextTimeout time.Duration
}

type SignupUsecase interface {
	Signup(*models.User) error
}

func NewSignupUsecase(userRepository repository.UserRepository, timeout time.Duration) SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		ContextTimeout: timeout,
	}
}

func (su *signupUsecase) Signup(user *models.User) error {
	_, err := su.userRepository.GetUserByEmail(user.Email)
	if err == nil {
		return utils.ErrEmailIsTaken
	}

	_, err = su.userRepository.GetUserByName(user.Name)
	if err == nil {
		return utils.ErrNameIsTaken
	}

	err = user.Password.Set(user.Password.Plaintext)
	if err != nil {
		return utils.ErrInternalServer
	}

	_, err = su.userRepository.CreateUser(user)
	if err != nil {
		return utils.ErrInternalServer
	}

	return nil
}
