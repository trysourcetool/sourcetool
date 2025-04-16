package main

import (
	"errors"
	"fmt"
	"time"
)

type TicketStatus string
type TicketPriority string

const (
	StatusOpen       TicketStatus = "open"
	StatusInProgress TicketStatus = "in_progress"
	StatusResolved   TicketStatus = "resolved"
	StatusClosed     TicketStatus = "closed"

	PriorityLow    TicketPriority = "low"
	PriorityMedium TicketPriority = "medium"
	PriorityHigh   TicketPriority = "high"
	PriorityUrgent TicketPriority = "urgent"
)

type Ticket struct {
	ID          string
	Title       string
	Description string
	CustomerID  string
	Status      TicketStatus
	Priority    TicketPriority
	Assignee    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

var (
	ErrTicketNil           = errors.New("ticket cannot be nil")
	ErrTitleRequired       = errors.New("title is required")
	ErrDescriptionRequired = errors.New("description is required")
	ErrCustomerIDRequired  = errors.New("customer ID is required")
	ErrStatusInvalid       = errors.New("invalid status")
	ErrPriorityInvalid     = errors.New("invalid priority")
)

func listTickets(title, customerID string, status TicketStatus, priority TicketPriority, assignee string) ([]Ticket, error) {
	tickets := generateTestTickets()

	if title == "" && customerID == "" && status == "" && priority == "" && assignee == "" {
		return tickets, nil
	}

	return filterTickets(tickets, title, customerID, status, priority, assignee), nil
}

func filterTickets(tickets []Ticket, title, customerID string, status TicketStatus, priority TicketPriority, assignee string) []Ticket {
	var filtered []Ticket
	for _, ticket := range tickets {
		if matchesFilter(ticket, title, customerID, status, priority, assignee) {
			filtered = append(filtered, ticket)
		}
	}
	return filtered
}

func matchesFilter(ticket Ticket, title, customerID string, status TicketStatus, priority TicketPriority, assignee string) bool {
	return (title == "" || ticket.Title == title) &&
		(customerID == "" || ticket.CustomerID == customerID) &&
		(status == "" || ticket.Status == status) &&
		(priority == "" || ticket.Priority == priority) &&
		(assignee == "" || ticket.Assignee == assignee)
}

func createTicket(t *Ticket) error {
	if t == nil {
		return ErrTicketNil
	}

	if err := validateTicket(t); err != nil {
		return err
	}

	if t.ID == "" {
		t.ID = fmt.Sprintf("ticket_%d", time.Now().UnixNano())
	}

	now := time.Now().UTC()
	if t.CreatedAt.IsZero() {
		t.CreatedAt = now
	}
	t.UpdatedAt = now

	return nil
}

func updateTicket(t *Ticket) error {
	if t == nil {
		return ErrTicketNil
	}

	if err := validateTicket(t); err != nil {
		return err
	}

	t.UpdatedAt = time.Now().UTC()
	return nil
}

func validateTicket(t *Ticket) error {
	if t.Title == "" {
		return ErrTitleRequired
	}
	if t.Description == "" {
		return ErrDescriptionRequired
	}
	if t.CustomerID == "" {
		return ErrCustomerIDRequired
	}
	if t.Status != "" && t.Status != StatusOpen && t.Status != StatusInProgress && t.Status != StatusResolved && t.Status != StatusClosed {
		return ErrStatusInvalid
	}
	if t.Priority != "" && t.Priority != PriorityLow && t.Priority != PriorityMedium && t.Priority != PriorityHigh && t.Priority != PriorityUrgent {
		return ErrPriorityInvalid
	}
	return nil
}

func generateTestTickets() []Ticket {
	baseTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return []Ticket{
		{
			ID:          "1",
			Title:       "Cannot login to account",
			Description: "User reported unable to login to their account",
			CustomerID:  "cust_001",
			Status:      StatusOpen,
			Priority:    PriorityHigh,
			Assignee:    "support_agent_1",
			CreatedAt:   baseTime.Add(24 * time.Hour * 0),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 0),
		},
		{
			ID:          "2",
			Title:       "Payment processing error",
			Description: "Payment is not being processed correctly",
			CustomerID:  "cust_002",
			Status:      StatusInProgress,
			Priority:    PriorityUrgent,
			Assignee:    "support_agent_2",
			CreatedAt:   baseTime.Add(24 * time.Hour * 1),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 1),
		},
		{
			ID:          "3",
			Title:       "Feature request: Dark mode",
			Description: "Customer would like to have dark mode option",
			CustomerID:  "cust_003",
			Status:      StatusOpen,
			Priority:    PriorityLow,
			Assignee:    "support_agent_1",
			CreatedAt:   baseTime.Add(24 * time.Hour * 2),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 2),
		},
		{
			ID:          "4",
			Title:       "Account deletion request",
			Description: "Customer wants to delete their account and all associated data",
			CustomerID:  "cust_004",
			Status:      StatusInProgress,
			Priority:    PriorityMedium,
			Assignee:    "support_agent_3",
			CreatedAt:   baseTime.Add(24 * time.Hour * 3),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 3),
		},
		{
			ID:          "5",
			Title:       "Subscription renewal failed",
			Description: "Automatic renewal failed due to expired credit card",
			CustomerID:  "cust_005",
			Status:      StatusResolved,
			Priority:    PriorityHigh,
			Assignee:    "support_agent_2",
			CreatedAt:   baseTime.Add(24 * time.Hour * 4),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 4),
		},
		{
			ID:          "6",
			Title:       "Mobile app crashes on startup",
			Description: "App crashes immediately after launching on iOS 17.2",
			CustomerID:  "cust_006",
			Status:      StatusOpen,
			Priority:    PriorityUrgent,
			Assignee:    "support_agent_4",
			CreatedAt:   baseTime.Add(24 * time.Hour * 5),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 5),
		},
		{
			ID:          "7",
			Title:       "Data export format issue",
			Description: "CSV export includes incorrect date format",
			CustomerID:  "cust_007",
			Status:      StatusInProgress,
			Priority:    PriorityMedium,
			Assignee:    "support_agent_1",
			CreatedAt:   baseTime.Add(24 * time.Hour * 6),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 6),
		},
		{
			ID:          "8",
			Title:       "API rate limiting concerns",
			Description: "Customer hitting rate limits during peak hours",
			CustomerID:  "cust_008",
			Status:      StatusOpen,
			Priority:    PriorityHigh,
			Assignee:    "support_agent_3",
			CreatedAt:   baseTime.Add(24 * time.Hour * 7),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 7),
		},
		{
			ID:          "9",
			Title:       "Billing address update",
			Description: "Need to update billing address for tax purposes",
			CustomerID:  "cust_009",
			Status:      StatusResolved,
			Priority:    PriorityLow,
			Assignee:    "support_agent_2",
			CreatedAt:   baseTime.Add(24 * time.Hour * 8),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 8),
		},
		{
			ID:          "10",
			Title:       "Integration with new CRM",
			Description: "Request for integration with Salesforce CRM",
			CustomerID:  "cust_010",
			Status:      StatusOpen,
			Priority:    PriorityMedium,
			Assignee:    "support_agent_4",
			CreatedAt:   baseTime.Add(24 * time.Hour * 9),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 9),
		},
		{
			ID:          "11",
			Title:       "Password reset not working",
			Description: "Password reset emails not being received",
			CustomerID:  "cust_011",
			Status:      StatusInProgress,
			Priority:    PriorityHigh,
			Assignee:    "support_agent_1",
			CreatedAt:   baseTime.Add(24 * time.Hour * 10),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 10),
		},
		{
			ID:          "12",
			Title:       "Report generation timeout",
			Description: "Large reports timing out after 5 minutes",
			CustomerID:  "cust_012",
			Status:      StatusOpen,
			Priority:    PriorityMedium,
			Assignee:    "support_agent_3",
			CreatedAt:   baseTime.Add(24 * time.Hour * 11),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 11),
		},
		{
			ID:          "13",
			Title:       "Two-factor authentication issues",
			Description: "2FA codes not being accepted",
			CustomerID:  "cust_013",
			Status:      StatusResolved,
			Priority:    PriorityHigh,
			Assignee:    "support_agent_2",
			CreatedAt:   baseTime.Add(24 * time.Hour * 12),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 12),
		},
		{
			ID:          "14",
			Title:       "Data import validation errors",
			Description: "CSV import failing with validation errors",
			CustomerID:  "cust_014",
			Status:      StatusOpen,
			Priority:    PriorityMedium,
			Assignee:    "support_agent_4",
			CreatedAt:   baseTime.Add(24 * time.Hour * 13),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 13),
		},
		{
			ID:          "15",
			Title:       "Email notifications delayed",
			Description: "System emails being delayed by 2-3 hours",
			CustomerID:  "cust_015",
			Status:      StatusInProgress,
			Priority:    PriorityHigh,
			Assignee:    "support_agent_1",
			CreatedAt:   baseTime.Add(24 * time.Hour * 14),
			UpdatedAt:   baseTime.Add(24 * time.Hour * 14),
		},
	}
}
