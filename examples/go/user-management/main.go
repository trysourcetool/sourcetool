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
	name := searchCols[0].TextInput("Name", textinput.Placeholder("Enter name to filter"))
	email := searchCols[1].TextInput("Email", textinput.Placeholder("Enter email to filter"))

	users, err := listUsers(name, email, 0, "", time.Time{})
	if err != nil {
		return err
	}

	baseCols := ui.Columns(2, columns.Weight(3, 1))
	table := baseCols[0].Table(
		users,
		table.Header("Users"),
		table.Height(10),
		table.ColumnOrder("ID", "Name", "Email", "Age", "Gender", "CreatedAt"),
		table.OnSelect(table.SelectionBehaviorRerun),
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

	form, submitted := baseCols[1].Form("Update", form.ClearOnSubmit(true))
	formName := form.TextInput("Name", textinput.Placeholder("Enter your name"), textinput.DefaultValue(defaultName), textinput.Required(true))
	formEmail := form.TextInput("Email", textinput.Placeholder("Enter your email"), textinput.DefaultValue(defaultEmail))
	formAge := form.NumberInput("Age", numberinput.MinValue(0), numberinput.MaxValue(100), numberinput.DefaultValue(float64(defaultAge)))
	formGender := form.Selectbox("Gender", selectbox.Options("male", "female"), selectbox.DefaultValue(defaultGender))

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

	form, submitted := ui.Form("Create User", form.ClearOnSubmit(true))
	formName := form.TextInput("Name", textinput.Placeholder("Enter user name"), textinput.Required(true))
	formEmail := form.TextInput("Email", textinput.Placeholder("Enter user email"))
	formAge := form.NumberInput("Age", numberinput.MinValue(0), numberinput.MaxValue(100))
	formGender := form.Selectbox("Gender", selectbox.Options("male", "female"))

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
		APIKey:   "your_development_api_key",
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
