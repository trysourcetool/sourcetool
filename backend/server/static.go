package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ServeStaticFiles configures the router to serve static files from the specified directory
func ServeStaticFiles(r chi.Router, staticDir string) {
	if staticDir == "" {
		staticDir = os.Getenv("STATIC_FILES_DIR")
		if staticDir == "" {
			staticDir = "./static" // Default to ./static if not specified
		}
	}

	// Create a file server handler
	fileServer := http.FileServer(http.Dir(staticDir))

	// First, check if index.html exists
	indexPath := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		fmt.Printf("WARNING: index.html not found at %s\n", indexPath)
		// Try to find it in parent directories
		if staticDir == "/app/static" {
			alternativePath := "/app/static-full/client/index.html"
			if _, err := os.Stat(alternativePath); !os.IsNotExist(err) {
				fmt.Printf("Found index.html at alternative location: %s\n", alternativePath)
				indexPath = alternativePath
			}
		}
	} else {
		fmt.Printf("Found index.html at: %s\n", indexPath)
	}

	// Serve static files
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// Get the path
		path := r.URL.Path

		// Check if the file exists
		filePath := filepath.Join(staticDir, path)
		_, err := os.Stat(filePath)
		fileExists := !os.IsNotExist(err)

		// Always handle API, WebSocket, and Swagger paths with the file server
		// or if the file actually exists on disk and is not index.html
		if strings.HasPrefix(path, "/api") ||
			strings.HasPrefix(path, "/ws") ||
			strings.HasPrefix(path, "/swagger") ||
			(fileExists && !strings.HasSuffix(filePath, "index.html")) {

			// For asset files, set appropriate cache headers
			if strings.HasPrefix(path, "/assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000")
			}

			fileServer.ServeHTTP(w, r)
			return
		}

		// For all other paths, serve index.html to support client-side routing
		// Set headers for SPA routing support
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		// Add headers to help with client-side routing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Log the path for debugging
		fmt.Printf("Serving index.html for client-side route: %s\n", path)

		// Serve the index.html file
		http.ServeFile(w, r, indexPath)
	})
}
