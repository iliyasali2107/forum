package service

import (
	"errors"
	"fmt"
	"time"

	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg/validator"

	"github.com/google/uuid"
)

var (
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidUsernameLen  = errors.New("username length out of range 32")
	ErrInvalidUsernameChar = errors.New("invalid username characters")
	ErrInternalServer      = errors.New("internal server error")
	ErrConfirmPassword     = errors.New("password doesn't match")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists          = errors.New("user already exists")
	ErrFormValidation      = errors.New("form validation failed")
)

type AuthService interface {
	Signup(*validator.Validator, *models.User) error
	Login(*models.User) error
	Logout(*models.User) error
	ParseToken(token string) (*models.User, error)
	DeleteToken(token string) error
}

type authService struct {
	ur repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		ur: userRepo,
	}
}

func (as *authService) Signup(v *validator.Validator, user *models.User) error {
	_, err := as.ur.GetUserByEmail(user.Email)
	if err == nil {
		return ErrUserExists
	}

	if ValidateUser(v, user); !v.Valid() {
		return ErrFormValidation
	}

	err = user.Password.Set(user.Password.Plaintext)
	if err != nil {
		return ErrInternalServer
	}
	//TODO: UNIQUE constraint failed: users.name
	_, err = as.ur.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		return ErrInternalServer
	}

	return nil
}

func (as *authService) Login(user *models.User) error {
	u, err := as.ur.GetUserByEmail(user.Email)
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
	err = as.ur.SaveToken(user)
	if err != nil {
		return err
	}

	return nil
}

func (as *authService) Logout(user *models.User) error {
	u, err := as.ur.GetUserByEmail(user.Email)
	if err != nil {
		return ErrUserNotFound
	}

	user.ID = u.ID

	return as.ur.DeleteToken(*user.Token)

}

func (as *authService) ParseToken(token string) (*models.User, error) {
	return as.ur.GetUserByToken(token)
}

func (as *authService) DeleteToken(token string) error {
	return as.ur.DeleteToken(token)
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be provided a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	fmt.Println(password)
	fmt.Println(len(password))
	v.Check(len(password) <= 72, "password", "must not be more than 500 bytes long")
}

func ValidateUser(v *validator.Validator, user *models.User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	ValidatePasswordPlaintext(v, user.Password.Plaintext)
}
