package service

import (
	"errors"

	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg/validator"
)

var (
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidUsernameLen  = errors.New("username length out of range 32")
	ErrInvalidUsernameChar = errors.New("invalid username characters")

	ErrConfirmPassword = errors.New("password doesn't match")
	ErrUserNotFound    = errors.New("user not found")
	ErrUserExist       = errors.New("user already exists")
)

type AuthService interface {
	CreateUser(*models.User) error
	GenerateToken(username, password string) (*models.User, error)
	ParseToken(token string) (*models.User, error)
	DeleteToken(token string) error
}

type authService struct {
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
}

func NewAuthService(authRepo repository.AuthRepository, userRepo repository.UserRepository) AuthService {
	return &authService{
		authRepository: authRepo,
		userRepository: userRepo,
	}
}

func (s *authService) CreateUser(user *models.User) error {
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

func (s *authService) Signup(v *validator.Validator, user *models.User) error {
	_, err := s.userRepository.GetUserByEmail(user.Email)
	if err == nil {
		return ErrUserExist
	}

	if ValidateUser(v, user); !v.Valid() {
		if _, ok := v.Errors["name"]; ok {
			return ErrInvalidUsername
		}

		if _, ok := v.Errors["email"]; ok {
			return ErrInvalidEmail
		}

		if _, ok := v.Errors["password"]; ok {
			return ErrInvalidPassword
		}
	}
	
	// TODO: ADD: user.Password.Hash

	_, err = s.authRepository.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be provided a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 500 bytes long")
}

func ValidateUser(v *validator.Validator, user *models.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.Plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
