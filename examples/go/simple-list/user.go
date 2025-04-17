package main

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Age       int
	Gender    string
	CreatedAt time.Time
}

func listUsers(name, email string) ([]User, error) {
	users := generateTestUsers()

	if name == "" && email == "" {
		return users, nil
	}

	return filterUsers(users, name, email), nil
}

func filterUsers(users []User, name, email string) []User {
	var filtered []User
	for _, user := range users {
		if matchesFilter(user, name, email) {
			filtered = append(filtered, user)
		}
	}
	return filtered
}

func matchesFilter(user User, name, email string) bool {
	return (name == "" || user.Name == name) &&
		(email == "" || user.Email == email)
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
