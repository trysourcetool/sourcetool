package main

import (
	"fmt"
	"log"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func listTicketsPage(ui sourcetool.UIBuilder) error {
	searchCols := ui.Columns(3)
	title := searchCols[0].TextInput("Title", textinput.WithPlaceholder("Enter title to filter"))
	customerID := searchCols[1].TextInput("Customer ID", textinput.WithPlaceholder("Enter customer ID to filter"))
	assignee := searchCols[2].TextInput("Assignee", textinput.WithPlaceholder("Enter assignee to filter"))

	statusCols := ui.Columns(2)
	status := statusCols[0].Selectbox("Status", selectbox.WithOptions(
		string(StatusOpen),
		string(StatusInProgress),
		string(StatusResolved),
		string(StatusClosed),
	))
	priority := statusCols[1].Selectbox("Priority", selectbox.WithOptions(
		string(PriorityLow),
		string(PriorityMedium),
		string(PriorityHigh),
		string(PriorityUrgent),
	))

	var statusValue string
	if status != nil {
		statusValue = status.Value
	}
	var priorityValue string
	if priority != nil {
		priorityValue = priority.Value
	}

	tickets, err := listTickets(title, customerID, TicketStatus(statusValue), TicketPriority(priorityValue), assignee)
	if err != nil {
		return err
	}

	baseCols := ui.Columns(2, columns.WithWeight(3, 1))
	ticketTable := baseCols[0].Table(
		tickets,
		table.WithHeight(10),
		table.WithColumnOrder("ID", "Title", "CustomerID", "Status", "Priority", "Assignee", "CreatedAt"),
		table.WithOnSelect(table.OnSelectRerun),
	)

	var defaultTitle, defaultDescription, defaultCustomerID, defaultAssignee string
	var defaultStatus TicketStatus
	var defaultPriority TicketPriority
	if ticketTable.Selection != nil && ticketTable.Selection.Row < len(tickets) {
		selectedData := tickets[ticketTable.Selection.Row]
		defaultTitle = selectedData.Title
		defaultDescription = selectedData.Description
		defaultCustomerID = selectedData.CustomerID
		defaultStatus = selectedData.Status
		defaultPriority = selectedData.Priority
		defaultAssignee = selectedData.Assignee
	}

	form, submitted := baseCols[1].Form("Update Ticket", form.WithClearOnSubmit(true))
	formTitle := form.TextInput("Title", textinput.WithPlaceholder("Enter ticket title"), textinput.WithDefaultValue(defaultTitle), textinput.WithRequired(true))
	formDescription := form.TextInput("Description", textinput.WithPlaceholder("Enter ticket description"), textinput.WithDefaultValue(defaultDescription), textinput.WithRequired(true))
	formCustomerID := form.TextInput("Customer ID", textinput.WithPlaceholder("Enter customer ID"), textinput.WithDefaultValue(defaultCustomerID), textinput.WithRequired(true))
	formStatus := form.Selectbox("Status", selectbox.WithOptions(
		string(StatusOpen),
		string(StatusInProgress),
		string(StatusResolved),
		string(StatusClosed),
	), selectbox.WithDefaultValue(string(defaultStatus)))
	formPriority := form.Selectbox("Priority", selectbox.WithOptions(
		string(PriorityLow),
		string(PriorityMedium),
		string(PriorityHigh),
		string(PriorityUrgent),
	), selectbox.WithDefaultValue(string(defaultPriority)))
	formAssignee := form.TextInput("Assignee", textinput.WithPlaceholder("Enter assignee"), textinput.WithDefaultValue(defaultAssignee))

	if submitted {
		// Use the form status and priority values directly
		formStatusValue := formStatus.Value
		formPriorityValue := formPriority.Value

		ticket := Ticket{
			Title:       formTitle,
			Description: formDescription,
			CustomerID:  formCustomerID,
			Status:      TicketStatus(formStatusValue),
			Priority:    TicketPriority(formPriorityValue),
			Assignee:    formAssignee,
		}
		if err := updateTicket(&ticket); err != nil {
			return err
		}
		ui.Markdown(fmt.Sprintf("Ticket updated successfully"))
	}

	return nil
}

func createTicketPage(ui sourcetool.UIBuilder) error {
	form, submitted := ui.Form("Create Ticket", form.WithClearOnSubmit(true))
	formTitle := form.TextInput("Title", textinput.WithPlaceholder("Enter ticket title"), textinput.WithRequired(true))
	formDescription := form.TextInput("Description", textinput.WithPlaceholder("Enter ticket description"), textinput.WithRequired(true))
	formCustomerID := form.TextInput("Customer ID", textinput.WithPlaceholder("Enter customer ID"), textinput.WithRequired(true))
	formStatus := form.Selectbox("Status", selectbox.WithOptions(
		string(StatusOpen),
		string(StatusInProgress),
		string(StatusResolved),
		string(StatusClosed),
	), selectbox.WithDefaultValue(string(StatusOpen)))
	formPriority := form.Selectbox("Priority", selectbox.WithOptions(
		string(PriorityLow),
		string(PriorityMedium),
		string(PriorityHigh),
		string(PriorityUrgent),
	), selectbox.WithDefaultValue(string(PriorityMedium)))
	formAssignee := form.TextInput("Assignee", textinput.WithPlaceholder("Enter assignee"))

	if submitted {
		// Use the form status and priority values directly
		formStatusValue := formStatus.Value
		formPriorityValue := formPriority.Value

		ticket := Ticket{
			Title:       formTitle,
			Description: formDescription,
			CustomerID:  formCustomerID,
			Status:      TicketStatus(formStatusValue),
			Priority:    TicketPriority(formPriorityValue),
			Assignee:    formAssignee,
		}
		if err := createTicket(&ticket); err != nil {
			return err
		}
		ui.Markdown(fmt.Sprintf("Ticket created successfully with ID: %s", ticket.ID))
	}

	return nil
}

func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "your_api_key",
		Endpoint: "ws://localhost:3000",
	})

	s.Page("/tickets", "Tickets", listTicketsPage)
	s.Page("/tickets/new", "Create ticket", createTicketPage)

	if err := s.Listen(); err != nil {
		log.Printf("Failed to listen sourcetool: %v", err)
		s.Close()
		return
	}
}
