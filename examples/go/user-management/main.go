package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/columns"
	"github.com/trysourcetool/sourcetool-go/form"
	"github.com/trysourcetool/sourcetool-go/numberinput"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	"github.com/trysourcetool/sourcetool-go/table"
	"github.com/trysourcetool/sourcetool-go/textinput"
	"golang.org/x/sync/errgroup"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Hello, World!")
}

func listUsersPage(ui sourcetool.UIBuilder) error {
	ui.Markdown("## Users")

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
		table.WithHeader("Users"),
		table.WithHeight(10),
		table.WithColumnOrder("ID", "Name", "Email", "Age", "Gender", "CreatedAt"),
		table.WithOnSelect(table.OnSelectRerun),
	)

	var defaultName, defaultEmail, defaultGender string
	var defaultAge int
	if table.Selection != nil {
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
	ui.Markdown("## Create New User")

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
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	fmt.Println("Starting server at http://localhost:8081/")

	// Replace with your own API key for development
	config := &sourcetool.Config{
		APIKey:   "development_1Tly8TW7bYn716Gzi7XPPjtL2jmdMA6G1Tly8TW7bYn716Gzi7X",
		Endpoint: "ws://localhost:3000",
	}
	s := sourcetool.New(config)

	var (
		groupAdmin           = "admin"
		groupUserAdmin       = "user_admin"
		groupCustomerSupport = "customer_support"
	)
	s.AccessGroups(groupAdmin)
	usersGroup := s.Group("/users")
	{
		usersGroup.AccessGroups(groupUserAdmin)
		usersGroup.Page("/", "Users", listUsersPage)

		csGroup := usersGroup.Group("/")
		{
			csGroup.AccessGroups(groupCustomerSupport)
			csGroup.Page("/new", "Create user", createUserPage)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Listen(); err != nil {
			return fmt.Errorf("failed to listen sourcetool: %v", err)
		}
		return nil
	})

	eg.Go(func() error {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("failed to start server: %v", err)
		}
		return nil
	})
	eg.Go(func() error {
		<-egCtx.Done()
		log.Println("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %v", err)
		}

		if err := s.Close(); err != nil {
			return fmt.Errorf("sourcetool shutdown failed: %v", err)
		}

		log.Println("Server shutdown complete")
		return nil
	})

	if err := eg.Wait(); err != nil {
		log.Printf("Error during shutdown: %v", err)
		os.Exit(1)
	}
}
