package main

import (
	"fmt"
	"time"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
}

type RefundRequest struct {
	UserID string
	Amount int
	Reason string
}

func listUsers(name, email string) ([]User, error) {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	users := []User{
		{"1", "John Doe", "john.doe@acme.com", baseTime},
		{"2", "Jane Smith", "jane.smith@acme.com", baseTime.Add(24 * time.Hour)},
		{"3", "Bob Lee", "bob.lee@acme.com", baseTime.Add(48 * time.Hour)},
	}
	var filtered []User
	for _, u := range users {
		if (name == "" || u.Name == name) && (email == "" || u.Email == email) {
			filtered = append(filtered, u)
		}
	}
	return filtered, nil
}

func refundStripe(req RefundRequest) error {
	fmt.Printf("Refund: userID=%s, amount=%d, reason=%s\n", req.UserID, req.Amount, req.Reason)
	return nil
}
