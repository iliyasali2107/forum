package usecase

import (
	"time"

	"forum/bootstrap"
	"forum/domain/models"
	"forum/domain/repository"

	"github.com/google/uuid"
)

type loginUsecase struct {
	userRepository repository.UserRepository
	Env            *bootstrap.Env
	ContextTimeout time.Duration
}

type LoginUsecase interface {
	Login(*models.User) error
}

func NewLoginUsecase(userRepository repository.UserRepository, env *bootstrap.Env, timeout time.Duration) LoginUsecase {
	return &loginUsecase{
		userRepository: userRepository,
		Env:            env,
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
	expiresAt := time.Now().Add(lu.Env.LoginExpireTime * time.Minute)
	user.Token = &sessionToken
	user.Expires = &expiresAt
	user.ID = u.ID
	err = lu.userRepository.SaveToken(user)
	if err != nil {
		return err
	}

	return nil
}
