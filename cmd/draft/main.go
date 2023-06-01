package main

import (
	"fmt"
	"forum/pkg/validator"
)

func main() {
	sl := []Comment{
		{
			id: 1,
		},
		{
			id: 2,
		},
	}
	for i := range sl {
		sl[i].id = 3
	}

	fmt.Println(sl)
}

func isValid(email string) bool {
	// emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	// return emailRegex.MatchString(email)
	return validator.EmailRX.MatchString(email)
}

func isValidName(name string) bool {
	return validator.NameRX.MatchString(name)
}

type Comment struct {
	id int
}
