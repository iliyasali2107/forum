package usecase

import (
	"forum/internal/models"
	"forum/internal/repository"
	"time"

	"github.com/google/uuid"
)

type loginUsecase struct {
	userRepository repository.UserRepository
}

type LoginUsecase interface {
	Login(*models.User) error
}

func NewLoginUsecase(userRepository repository.UserRepository, timeout time.Duration) LoginUsecase {
	return &loginUsecase{
		userRepository: userRepository,
	}
}

func (lu *loginUsecase) Login(user *models.User) error {
	u, err := lu.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		return ErrUserNotFound
	}

	ok, err := u.Password.Matches(user.Password.Plaintext)
	if err != nil || !ok {
		return err
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)
	user.Token = &sessionToken
	user.Expires = &expiresAt
	user.ID = u.ID
	err = lu.userRepository.SaveToken(user)
	if err != nil {
		return err
	}

	return nil
}
