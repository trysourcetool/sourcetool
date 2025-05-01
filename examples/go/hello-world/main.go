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
	"github.com/trysourcetool/sourcetool-go/button"
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

func helloPage(ui sourcetool.UIBuilder) error {
	ui.Markdown("# Hello, Sourcetool!")
	ui.Markdown("This is a simple example demonstrating the basic usage of the Sourcetool Go SDK.")

	name := ui.TextInput("Your Name", textinput.WithPlaceholder("Enter your name"))

	if name != "" {
		ui.Markdown(fmt.Sprintf("## Hello, %s!", name))
		ui.Markdown("Welcome to Sourcetool!")
	}

	clicked := ui.Button("Say Hello", button.WithDisabled(false))
	if clicked {
		ui.Markdown("👋 Hello from the button click!")
	}

	return nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", helloHandler)

	server := &http.Server{
		Addr:    ":8082",
		Handler: mux,
	}

	fmt.Println("Starting server at http://localhost:8082/")

	// Replace with your own API key for development
	config := &sourcetool.Config{
		APIKey:   "your_development_api_key",
		Endpoint: "ws://localhost:3000",
	}
	s := sourcetool.New(config)

	s.Page("/hello", "Hello", helloPage)

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
