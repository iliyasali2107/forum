package service

import (
	"errors"
	"fmt"
	"net/mail"

	"forum/internal/models"
	"forum/internal/repository"
)

var (
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidUsernameLen  = errors.New("username length out of range 32")
	ErrInvalidUsernameChar = errors.New("invalid username characters")
	ErrConfirmPassword     = errors.New("password doesn't match")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExist           = errors.New("user already exists")
)

type AuthService interface {
	CreateUser(*models.User) error
	GenerateToken(username, password string) (*models.User, error)
	ParseToken(token string) (*models.User, error)
	DeleteToken(token string) error
}

type authService struct {
	Repository repository.AuthRepository
}

func NewAuthService(repository repository.AuthRepository) AuthService {
	return &authService{
		Repository: repository,
	}
}

func (s *authService) CreateUser(user *models.User) error {
	if _, err := s.Repository.GetUser(user.Name); err == nil {
		return fmt.Errorf("service: CreateUser: get user: %W", err)
	}

	return nil
}

func (s *authService) GenerateToken(username, password string) (*models.User, error) {
	return nil, nil
}

func (s *authService) ParseToken(token string) (*models.User, error) {
	return nil, nil
}

func (s *authService) DeleteToken(token string) error {
	return nil
}

func checkUser(user models.User) error {
	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrInvalidEmail)
	}

	for _, char := range user.Name {
		if char < 32 || char > 126 {
			return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrInvalidUsernameChar)
		}
	}

	if len(user.Name) < 1 || len(user.Name) >= 36 {
		return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrInvalidUsernameLen)
	}

	// if user.Password != user.ConfirmPassword {
	// 	return fmt.Errorf("service: CreateUser: checkUser err: %w", ErrConfirmPassword)
	// }
	// return nil
	return nil
}
