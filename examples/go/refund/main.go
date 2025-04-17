package main

import (
	"fmt"
	"log"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/numberinput"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func refundPage(ui sourcetool.UIBuilder) error {
	searchCols := ui.Columns(2)
	name := searchCols[0].TextInput("Name", textinput.WithPlaceholder("Enter user name to filter"))
	email := searchCols[1].TextInput("Email", textinput.WithPlaceholder("Enter email to filter"))

	users, err := listUsers(name, email)
	if err != nil {
		return err
	}

	tableWidget := ui.Table(
		users,
		table.WithHeight(10),
		table.WithColumnOrder("ID", "Name", "Email", "CreatedAt"),
		table.WithOnSelect(table.OnSelectRerun),
	)

	var selectedUser *User
	if tableWidget.Selection != nil && tableWidget.Selection.Row < len(users) {
		selectedUser = &users[tableWidget.Selection.Row]
	}

	if selectedUser != nil {
		formWidget, submitted := ui.Form("Refund", form.WithClearOnSubmit(true))
		amount := formWidget.NumberInput("Amount", numberinput.WithMinValue(1), numberinput.WithRequired(true))
		reason := formWidget.TextInput("Reason", textinput.WithPlaceholder("Enter refund reason"), textinput.WithRequired(true))

		if submitted {
			refundReq := RefundRequest{
				UserID: selectedUser.ID,
				Amount: int(*amount),
				Reason: reason,
			}
			if err := refundStripe(refundReq); err != nil {
				return err
			}
			ui.Markdown(fmt.Sprintf("Refund processed for user %s (%s), amount: %d", selectedUser.Name, selectedUser.Email, refundReq.Amount))
		}
	}

	return nil
}

func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "your_api_key",
		Endpoint: "ws://localhost:3000",
	})

	s.Page("/refunds", "Refunds", refundPage)

	if err := s.Listen(); err != nil {
		log.Printf("Failed to listen sourcetool: %v", err)
		s.Close()
		return
	}
}
