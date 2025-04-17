package main

import (
	"log"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func listUsersPage(ui sourcetool.UIBuilder) error {
	searchCols := ui.Columns(2)
	name := searchCols[0].TextInput("Name", textinput.WithPlaceholder("Enter name to filter"))
	email := searchCols[1].TextInput("Email", textinput.WithPlaceholder("Enter email to filter"))

	users, err := listUsers(name, email)
	if err != nil {
		return err
	}

	ui.Table(
		users,
		table.WithHeight(10),
		table.WithColumnOrder("ID", "Name", "Email", "Age", "Gender", "CreatedAt"),
	)

	return nil
}

func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "your_api_key",
		Endpoint: "ws://localhost:3000",
	})

	s.Page("/users", "Users", listUsersPage)

	if err := s.Listen(); err != nil {
		log.Printf("Failed to listen sourcetool: %v", err)
		s.Close()
		return
	}
}
