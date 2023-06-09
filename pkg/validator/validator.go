package validator

import (
	"forum/domain/models"
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	NameRX  = regexp.MustCompile(`^[a-zA-Z0-9]{3,50}$`)
)

func CreatePostValidation(post *models.Post) map[string]string {
	errors := make(map[string]string)
	if post.Title == "" {
		errors["title"] = "must be provided"
	}

	if len(post.Title) > 100 {
		errors["title"] = "must not be more than 100 chars"
	}

	if post.Content == "" {
		errors["content"] = "must be provided"
	}

	if len(post.Content) > 100 {
		errors["content"] = "must not be more than 100 chars"
	}

	return errors
}

func SignupValidation(user *models.User) map[string]string {
	errors := make(map[string]string)

	if user.Name == "" {
		errors["name"] = "must be provided"
	}

	if !NameRX.MatchString(user.Name) {
		errors["name"] = "must be 3 to 50 chars and have alphanumeric chars"
	}

	if !EmailRX.MatchString(user.Email) {
		errors["email"] = "must be correctly formatted"
	}

	if !isValidPassword(user.Password.Plaintext) {
		errors["password"] = "at least 1 lower, 1 digit & 6 to 50 chars"
	}

	return errors
}

func isValidPassword(pass string) bool {
	tests := []string{".{6,}", "[a-z]", "[0-9]"}
	for _, test := range tests {
		valid, _ := regexp.MatchString(test, pass)
		if !valid {
			return false
		}
	}
	return true
}
