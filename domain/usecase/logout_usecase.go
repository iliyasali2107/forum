package usecase

import (
	"fmt"
	"time"

	"forum/domain/models"
	"forum/domain/repository"
)

type logoutUsecase struct {
	userRepository repository.UserRepository
	ContextTimeout time.Duration
}
type LogoutUsecase interface {
	Logout(user *models.User) error
	DeleteToken(token string) error
}

func NewLogoutUsecase(userRepository repository.UserRepository, timeout time.Duration) LogoutUsecase {
	return &logoutUsecase{
		userRepository: userRepository,
		ContextTimeout: timeout,
	}
}

func (lu *logoutUsecase) Logout(user *models.User) error {
	u, err := lu.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.ID = u.ID

	return nil

	// return lu.userRepository.DeleteToken(*user.Token)
}

func (lu *logoutUsecase) ParseToken(token string) (*models.User, error) {
	return lu.userRepository.GetUserByToken(token)
}

func (lu *logoutUsecase) DeleteToken(token string) error {
	return lu.userRepository.DeleteToken(token)
}
