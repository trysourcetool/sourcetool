package main

import (
	"fmt"
	"log"
	"time"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/numberinput"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

func listUsersPage(ui sourcetool.UIBuilder) error {
	searchCols := ui.Columns(2)
	name := searchCols[0].TextInput("Name", textinput.WithPlaceholder("Enter name to filter"))
	email := searchCols[1].TextInput("Email", textinput.WithPlaceholder("Enter email to filter"))

	users, err := listUsers(name, email, 0, "", time.Time{})
	if err != nil {
		return err
	}

	baseCols := ui.Columns(2, columns.WithWeight(3, 1))
	table := baseCols[0].Table(
		users,
		table.WithHeight(10),
		table.WithColumnOrder("ID", "Name", "Email", "Age", "Gender", "CreatedAt"),
		table.WithOnSelect(table.OnSelectRerun),
	)

	var defaultName, defaultEmail, defaultGender string
	var defaultAge int
	if table.Selection != nil && table.Selection.Row < len(users) {
		selectedData := users[table.Selection.Row]
		defaultName = selectedData.Name
		defaultEmail = selectedData.Email
		defaultAge = selectedData.Age
		defaultGender = selectedData.Gender
	}

	form, submitted := baseCols[1].Form("Update", form.WithClearOnSubmit(true))
	formName := form.TextInput("Name", textinput.WithPlaceholder("Enter your name"), textinput.WithDefaultValue(defaultName), textinput.WithRequired(true))
	formEmail := form.TextInput("Email", textinput.WithPlaceholder("Enter your email"), textinput.WithDefaultValue(defaultEmail))
	formAge := form.NumberInput("Age", numberinput.WithMinValue(0), numberinput.WithMaxValue(100), numberinput.WithDefaultValue(float64(defaultAge)))
	formGender := form.Selectbox("Gender", selectbox.WithOptions("male", "female"), selectbox.WithDefaultValue(defaultGender))

	if submitted {
		user := User{
			Name:   formName,
			Email:  formEmail,
			Age:    int(*formAge),
			Gender: formGender.Value,
		}
		if err := createUser(&user); err != nil {
			return err
		}
	}

	return nil
}

func createUserPage(ui sourcetool.UIBuilder) error {
	form, submitted := ui.Form("Create User", form.WithClearOnSubmit(true))
	formName := form.TextInput("Name", textinput.WithPlaceholder("Enter user name"), textinput.WithRequired(true))
	formEmail := form.TextInput("Email", textinput.WithPlaceholder("Enter user email"))
	formAge := form.NumberInput("Age", numberinput.WithMinValue(0), numberinput.WithMaxValue(100))
	formGender := form.Selectbox("Gender", selectbox.WithOptions("male", "female"))

	if submitted {
		user := User{
			Name:   formName,
			Email:  formEmail,
			Age:    int(*formAge),
			Gender: formGender.Value,
		}
		if err := createUser(&user); err != nil {
			return err
		}
		ui.Markdown(fmt.Sprintf("User created successfully with ID: %s", user.ID))
	}

	return nil
}

func main() {
	s := sourcetool.New(&sourcetool.Config{
		APIKey:   "development_MEst0aYqiDmlAgS1laq5tXFA94ERtB8rMEst0aYqiDmlAgS1laq",
		Endpoint: "ws://localhost:3000",
	})

	s.Page("/users", "Users", listUsersPage)
	s.Page("/users/new", "Create user", createUserPage)

	if err := s.Listen(); err != nil {
		log.Printf("Failed to listen sourcetool: %v", err)
		s.Close()
		return
	}
}
