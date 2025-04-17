package main

import (
	"fmt"
	"log"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/numberinput"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	"github.com/trysourcetool/sourcetool-go/textinput"
)

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
		APIKey:   "your_api_key",
		Endpoint: "ws://localhost:3000",
	})

	s.Page("/users/new", "Create user", createUserPage)

	if err := s.Listen(); err != nil {
		log.Printf("Failed to listen sourcetool: %v", err)
		s.Close()
		return
	}
}
