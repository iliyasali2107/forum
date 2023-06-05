package utils

import "errors"

var (
	ErrCategoryNotFound    = errors.New("category not found")
	ErrNameIsTaken         = errors.New("name is taken")
	ErrEmailIsTaken        = errors.New("email is taken")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidUsername     = errors.New("invalid username")
	ErrInvalidUsernameLen  = errors.New("username length out of range 32")
	ErrInvalidUsernameChar = errors.New("invalid username characters")
	ErrInternalServer      = errors.New("domain server error")
	ErrConfirmPassword     = errors.New("password doesn't match")
	ErrUserNotFound        = errors.New("user not found")
	ErrUserExists          = errors.New("user already exists")
	ErrFormValidation      = errors.New("form validation failed")
	ErrNoPosts             = errors.New("there is no posts")
)
