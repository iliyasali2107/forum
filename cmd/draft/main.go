package main

import (
	"forum/pkg/validator"
)

func main() {
}

func isValidName(name string) bool {
	return validator.NameRX.MatchString(name)
}

type Comment struct {
	id int
}
