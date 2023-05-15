package usecase

import (
	"forum/internal/models"
	"forum/internal/repository"
)

type logoutUsecase struct {
	userRepository repository.UserRepository
}
type LogoutUsecase interface {
	Logout(user *models.User) error
}

func NewLogoutUsecase(userRepository repository.UserRepository) LogoutUsecase {
	return &logoutUsecase{
		userRepository: userRepository,
	}
}

func (lu *logoutUsecase) Logout(user *models.User) error {
	u, err := lu.userRepository.GetUserByEmail(user.Email)
	if err != nil {
		return ErrUserNotFound
	}

	user.ID = u.ID

	return lu.userRepository.DeleteToken(*user.Token)
}

func (lu *logoutUsecase) ParseToken(token string) (*models.User, error) {
	return lu.userRepository.GetUserByToken(token)
}

func (lu *logoutUsecase) DeleteToken(token string) error {
	return lu.userRepository.DeleteToken(token)
}
