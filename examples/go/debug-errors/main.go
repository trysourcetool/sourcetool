package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/trysourcetool/sourcetool-go"
	"github.com/trysourcetool/sourcetool-go/button"
	"github.com/trysourcetool/sourcetool-go/numberinput"
	"github.com/trysourcetool/sourcetool-go/selectbox"
	"github.com/trysourcetool/sourcetool-go/textinput"
	"golang.org/x/sync/errgroup"
)

// Custom error types for demonstration
type ValidationError struct {
	Field string
	Msg   string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Validation error on field '%s': %s", e.Field, e.Msg)
}

type DatabaseError struct {
	Operation string
	Err       error
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("Database error during %s: %v", e.Operation, e.Err)
}

func (e DatabaseError) Unwrap() error {
	return e.Err
}

// Simulated database operations that will fail
func simulateDBQuery() error {
	return DatabaseError{
		Operation: "query",
		Err:       errors.New("connection timeout"),
	}
}

func simulateDBInsert() error {
	return DatabaseError{
		Operation: "insert",
		Err:       errors.New("duplicate key violation"),
	}
}

func simulateDBUpdate() error {
	return DatabaseError{
		Operation: "update",
		Err:       errors.New("record not found"),
	}
}

// Page that demonstrates various error scenarios
func errorDemoPage(ui sourcetool.UIBuilder) error {
	ui.Markdown("# Error Demonstration Page")
	ui.Markdown("This page demonstrates various error scenarios in Sourcetool.")
	ui.Markdown("**Note:** This page intentionally returns errors to demonstrate error handling in Sourcetool.")

	// Section for returning errors from the page function
	ui.Markdown("## Return Errors")
	ui.Markdown("Select an error type to return from the page function:")
	ui.Markdown("When an error is returned, Sourcetool will handle it and display an error message.")

	errorType := ui.Selectbox(
		"Error Type",
		selectbox.WithOptions(
			"None",
			"Validation Error",
			"Database Query Error",
			"Database Insert Error",
			"Database Update Error",
			"Panic",
			"Nil Pointer",
			"Index Out of Range",
			"Division by Zero",
		),
	)

	triggerError := ui.Button("Return Error", button.WithDisabled(false))

	if triggerError {
		switch errorType.Value {
		case "Validation Error":
			return ValidationError{
				Field: "username",
				Msg:   "Username cannot be empty",
			}
		case "Database Query Error":
			return simulateDBQuery()
		case "Database Insert Error":
			return simulateDBInsert()
		case "Database Update Error":
			return simulateDBUpdate()
		case "Panic":
			panic("Intentional panic for debugging purposes")
		case "Nil Pointer":
			var p *int
			// This will cause a nil pointer dereference
			_ = *p
		case "Index Out of Range":
			arr := []int{1, 2, 3}
			// This will cause an index out of range error
			_ = arr[10]
		case "Division by Zero":
			// This will cause a division by zero error
			var zero int
			_ = 1 / zero // This will panic at runtime, not compile time
		}
	}

	// Form with validation errors
	ui.Markdown("## Form with Validation Errors")
	ui.Markdown("Submit this form to return validation errors from the page function:")

	form, submitted := ui.Form("Submit Form")

	username := form.TextInput("Username", textinput.WithPlaceholder("Enter username"))
	email := form.TextInput("Email", textinput.WithPlaceholder("Enter email"))
	age := form.NumberInput("Age", numberinput.WithMinValue(0), numberinput.WithMaxValue(120))

	if submitted {
		// Validate form inputs
		if username == "" {
			return ValidationError{
				Field: "username",
				Msg:   "Username cannot be empty",
			}
		}

		if email == "" || !strings.Contains(email, "@") {
			return ValidationError{
				Field: "email",
				Msg:   "Invalid email format",
			}
		}

		if age == nil || *age < 18 {
			return ValidationError{
				Field: "age",
				Msg:   "Age must be at least 18",
			}
		}

		// Simulate a database error on form submission
		return simulateDBInsert()
	}

	// Demonstrate error handling with recovery
	ui.Markdown("## Error Recovery Demo")

	recoverDemo := ui.Button("Test Recovery", button.WithDisabled(false))

	if recoverDemo {
		func() {
			defer func() {
				if r := recover(); r != nil {
					ui.Markdown(fmt.Sprintf("**Recovered from panic:** %v", r))
				}
			}()

			// This will panic but be recovered
			panic("This panic will be recovered")
		}()
	}

	return nil
}

// Error handler middleware
func errorHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered in middleware: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	// Add error handler middleware
	handler := errorHandlerMiddleware(mux)

	server := &http.Server{
		Addr:    ":8083",
		Handler: handler,
	}

	fmt.Println("Starting server at http://localhost:8083/")

	// Replace with your own API key for development
	config := &sourcetool.Config{
		APIKey:   "your_development_api_key",
		Endpoint: "ws://localhost:3000",
	}
	s := sourcetool.New(config)

	// Register the error demo page
	s.AccessGroups("admin").Page("/error-demo", "Error Demo", errorDemoPage)

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
