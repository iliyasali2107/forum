package service

import (
	"errors"
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
)

type AuthService interface {
	Signup(*validator.Validator, *models.User) []error
	Login(*validator.Validator, *models.User) error
	Logout(*models.User) error
	GenerateToken(username, password string) (*models.User, error)
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

func (as *authService) GenerateToken(username, password string) (*models.User, error) {
	return nil, nil
}

func (as *authService) ParseToken(token string) (*models.User, error) {
	return nil, nil
}

func (as *authService) DeleteToken(token string) error {
	return nil
}

func (as *authService) Signup(v *validator.Validator, user *models.User) []error {
	errs := []error{}
	_, err := as.ur.GetUserByEmail(user.Email)
	if err == nil {
		errs = append(errs, ErrUserExists)
		return errs
	}

	if ValidateUser(v, user); !v.Valid() {
		if _, ok := v.Errors["name"]; ok {
			errs = append(errs, ErrInvalidUsername)
		}

		if _, ok := v.Errors["email"]; ok {
			errs = append(errs, ErrInvalidEmail)
		}

		if _, ok := v.Errors["password"]; ok {
			errs = append(errs, ErrInvalidPassword)
		}
	}

	err = user.Password.Set(user.Password.Plaintext)
	if err != nil {
		errs = append(errs, ErrInternalServer) //TODO: should return http.InternalServerError OR err (may be)
		return errs
	}

	_, err = as.ur.CreateUser(user)
	if err != nil {
		errs = append(errs, ErrInternalServer)
		return errs
	}

	return nil
}

func (as *authService) Login(v *validator.Validator, user *models.User) error {
	u, err := as.ur.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}

	ok, err := u.Password.Matches(user.Password.Plaintext)
	if err != nil || !ok {
		return ErrInvalidPassword
	}

	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(30 * time.Minute)

	user.Token = sessionToken
	user.Expires = expiresAt
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

	return as.ur.DeleteToken(user.ID)
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

	if user.Password.Plaintext != "" {
		ValidatePasswordPlaintext(v, user.Password.Plaintext)
	}

	if user.Password.Hash == nil {
		panic("missing password hash for user")
	}
}
