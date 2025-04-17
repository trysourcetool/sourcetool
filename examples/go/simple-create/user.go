package main

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	Age       int
	Gender    string
	CreatedAt time.Time
}

var (
	ErrUserNil        = errors.New("user cannot be nil")
	ErrNameRequired   = errors.New("name is required")
	ErrEmailRequired  = errors.New("email is required")
	ErrAgeInvalid     = errors.New("age must be positive")
	ErrGenderRequired = errors.New("gender is required")
	ErrGenderInvalid  = errors.New("gender must be either 'male' or 'female'")
)

func createUser(u *User) error {
	if u == nil {
		return ErrUserNil
	}

	if err := validateUser(u); err != nil {
		return err
	}

	if u.ID == "" {
		u.ID = fmt.Sprintf("user_%d", time.Now().UnixNano())
	}

	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}

	return nil
}

func validateUser(u *User) error {
	if u.Name == "" {
		return ErrNameRequired
	}
	if u.Email == "" {
		return ErrEmailRequired
	}
	if u.Age <= 0 {
		return ErrAgeInvalid
	}
	if u.Gender == "" {
		return ErrGenderRequired
	}
	if u.Gender != "male" && u.Gender != "female" {
		return ErrGenderInvalid
	}
	return nil
}
