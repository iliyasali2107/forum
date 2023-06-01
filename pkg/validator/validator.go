package validator

import (
	"regexp"

	"forum/domain/models"
)

var (
	EmailRX    = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	NameRX     = regexp.MustCompile(`^[a-zA-Z0-9]{3,50}$`)
	PasswordRX = regexp.MustCompile(`[a-z]+[A-Z]+[0-9]+[@$!%*?&]+[A-Za-z\d@$!%*?&]{8,46}`)
)

// emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Validator Define a new Validator type which contains a map of validation errors.
type Validator struct {
	Errors map[string]string
}

// NewValidator New is a helper which creates a new Validator isntance with an empty errors map.
func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid returns true if errors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message to the map (so long as no entry already exists for the given key).
func (v *Validator) AddError(key, message string) {
	v.Errors[key] = message
}

// Check adds an error message to the map only if a validation check is not 'ok'
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// In returns true if a specific value in a list of strings.
func In(value string, list ...string) bool {
	for i := range list {
		if value == list[i] {
			return true
		}
	}
	return false
}

// Matches returns true if a string value matches a specific regexp pattern
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique returns true if all string values in a slice are unique
func Unique(values []string) bool {
	uniqueValues := make(map[string]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

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

	if !PasswordRX.MatchString(user.Password.Plaintext) {
		errors["password"] = "at least 1 upper, 1 lower, 1 digit, 1 special char & 8 to 50 chars"
	}

	return errors
}
