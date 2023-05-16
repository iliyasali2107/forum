package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password Password
	Token    *string
	Expires  *time.Time
}

type Password struct {
	Plaintext string
	Hash      []byte
}

func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.Plaintext = plaintextPassword
	p.Hash = hash

	return nil
}

func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		return false, err
	}

	return true, nil
}
