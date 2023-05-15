package usecase

import (
	"forum/internal/models"
	"forum/internal/repository"
	"forum/pkg/validator"
)

type signupUsecase struct {
	userRepository repository.UserRepository
}

type SignupUsecase interface {
	Signup(*validator.Validator, *models.User) error
}

func NewSignupUsecase(userRepository repository.UserRepository) SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
	}
}

func (su *signupUsecase) Signup(v *validator.Validator, user *models.User) error {
	_, err := su.userRepository.GetUserByEmail(user.Email)
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
	// TODO: UNIQUE constraint failed: users.name
	_, err = su.userRepository.CreateUser(user)
	if err != nil {
		return ErrInternalServer
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

	ValidatePasswordPlaintext(v, user.Password.Plaintext)
}
