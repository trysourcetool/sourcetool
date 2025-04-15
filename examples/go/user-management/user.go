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

func listUsers(name, email string, age int, gender string, createdAt time.Time) ([]User, error) {
	users := generateTestUsers()

	if name == "" && email == "" && age == 0 && gender == "" && createdAt.IsZero() {
		return users, nil
	}

	return filterUsers(users, name, email, age, gender, createdAt), nil
}

func filterUsers(users []User, name, email string, age int, gender string, createdAt time.Time) []User {
	var filtered []User
	for _, user := range users {
		if matchesFilter(user, name, email, age, gender, createdAt) {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func matchesFilter(user User, name, email string, age int, gender string, createdAt time.Time) bool {
	return (name == "" || user.Name == name) &&
		(email == "" || user.Email == email) &&
		(age == 0 || user.Age == age) &&
		(gender == "" || user.Gender == gender) &&
		(createdAt.IsZero() || user.CreatedAt.Equal(createdAt))
}

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

func generateTestUsers() []User {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []User{
		{ID: "1", Name: "John Doe 001", Email: "john.doe+001@acme.com", Age: 25, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 0)},
		{ID: "2", Name: "John Doe 002", Email: "john.doe+002@acme.com", Age: 30, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 1)},
		{ID: "3", Name: "Jane Doe 003", Email: "jane.doe+003@acme.com", Age: 35, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 2)},
		{ID: "4", Name: "John Doe 004", Email: "john.doe+004@acme.com", Age: 28, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 3)},
		{ID: "5", Name: "Jane Doe 005", Email: "jane.doe+005@acme.com", Age: 32, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 4)},
		{ID: "6", Name: "John Doe 006", Email: "john.doe+006@acme.com", Age: 27, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 5)},
		{ID: "7", Name: "Jane Doe 007", Email: "jane.doe+007@acme.com", Age: 31, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 6)},
		{ID: "8", Name: "John Doe 008", Email: "john.doe+008@acme.com", Age: 29, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 7)},
		{ID: "9", Name: "Jane Doe 009", Email: "jane.doe+009@acme.com", Age: 33, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 8)},
		{ID: "10", Name: "John Doe 010", Email: "john.doe+010@acme.com", Age: 26, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 9)},
		{ID: "11", Name: "Jane Doe 011", Email: "jane.doe+011@acme.com", Age: 34, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 10)},
		{ID: "12", Name: "John Doe 012", Email: "john.doe+012@acme.com", Age: 28, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 11)},
		{ID: "13", Name: "Jane Doe 013", Email: "jane.doe+013@acme.com", Age: 30, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 12)},
		{ID: "14", Name: "John Doe 014", Email: "john.doe+014@acme.com", Age: 32, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 13)},
		{ID: "15", Name: "Jane Doe 015", Email: "jane.doe+015@acme.com", Age: 29, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 14)},
		{ID: "16", Name: "John Doe 016", Email: "john.doe+016@acme.com", Age: 31, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 15)},
		{ID: "17", Name: "Jane Doe 017", Email: "jane.doe+017@acme.com", Age: 27, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 16)},
		{ID: "18", Name: "John Doe 018", Email: "john.doe+018@acme.com", Age: 33, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 17)},
		{ID: "19", Name: "Jane Doe 019", Email: "jane.doe+019@acme.com", Age: 35, Gender: "female", CreatedAt: baseTime.Add(24 * time.Hour * 18)},
		{ID: "20", Name: "John Doe 020", Email: "john.doe+020@acme.com", Age: 30, Gender: "male", CreatedAt: baseTime.Add(24 * time.Hour * 19)},
	}
}
