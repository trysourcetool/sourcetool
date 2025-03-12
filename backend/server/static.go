package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

// ServeStaticFiles configures the router to serve static files from the specified directory.
func ServeStaticFiles(r chi.Router, staticDir string) {
	staticDir = getStaticDir(staticDir)
	if staticDir == "" {
		fmt.Println("Static file serving is disabled (local environment)")
		return
	}
	fileServer := http.FileServer(http.Dir(staticDir))
	indexPath := findIndexFile(staticDir)

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		filePath := filepath.Join(staticDir, path)

		if shouldServeFile(path, filePath) {
			serveFile(w, r, fileServer, path)
			return
		}

		serveIndexFile(w, r, indexPath)
	})
}

func getStaticDir(staticDir string) string {
	if staticDir == "" {
		staticDir = os.Getenv("STATIC_FILES_DIR")
		if staticDir == "" {
			// In local environment, return empty string to disable static serving
			return ""
		}
	}
	return staticDir
}

func findIndexFile(staticDir string) string {
	indexPath := filepath.Join(staticDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		if staticDir == "/app/static" {
			alternativePath := "/app/static-full/client/index.html"
			if _, err := os.Stat(alternativePath); !os.IsNotExist(err) {
				fmt.Printf("Found index.html at alternative location: %s\n", alternativePath)
				return alternativePath
			}
		}
		fmt.Printf("WARNING: index.html not found at %s\n", indexPath)
	} else {
		fmt.Printf("Found index.html at: %s\n", indexPath)
	}
	return indexPath
}

func shouldServeFile(path, filePath string) bool {
	// Check if the file exists
	_, err := os.Stat(filePath)
	fileExists := !os.IsNotExist(err)

	// Always handle API, WebSocket, and Swagger paths with the file server
	// or if the file actually exists on disk and is not index.html
	return strings.HasPrefix(path, "/api/") ||
		strings.HasPrefix(path, "/ws") ||
		strings.HasPrefix(path, "/swagger") ||
		(fileExists && !strings.HasSuffix(filePath, "index.html"))
}

func serveFile(w http.ResponseWriter, r *http.Request, fileServer http.Handler, path string) {
	if strings.HasPrefix(path, "/assets/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	fileServer.ServeHTTP(w, r)
}

func serveIndexFile(w http.ResponseWriter, r *http.Request, indexPath string) {
	// Set headers for SPA routing support
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")

	fmt.Printf("Serving index.html for client-side route: %s\n", r.URL.Path)
	http.ServeFile(w, r, indexPath)
}
