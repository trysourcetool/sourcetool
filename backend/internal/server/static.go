package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
)

func serveStaticFiles(r chi.Router, staticDir string) {
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
			alternativePath := "/app/static-full/index.html"
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

func setSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

	if config.Config.Env == config.EnvLocal {
		return
	}

	if config.Config.IsCloudEdition {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self' *.trysourcetool.com; "+
				"script-src 'self' 'unsafe-inline' 'unsafe-eval' *.trysourcetool.com; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src * data: blob:; "+
				"font-src * data:; "+
				"media-src *; "+
				"connect-src 'self' wss: *.trysourcetool.com https:;")
	} else {
		w.Header().Set("Content-Security-Policy",
			"default-src * 'unsafe-inline' 'unsafe-eval'; "+
				"img-src * data: blob:; "+
				"font-src * data:; "+
				"connect-src * ws: wss:;")
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, fileServer http.Handler, path string) {
	setSecurityHeaders(w)
	if strings.HasPrefix(path, "/assets/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}
	fileServer.ServeHTTP(w, r)
}

func serveIndexFile(w http.ResponseWriter, r *http.Request, indexPath string) {
	// Set headers for SPA routing support
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	setSecurityHeaders(w)

	fmt.Printf("Serving index.html for client-side route: %s\n", r.URL.Path)
	http.ServeFile(w, r, indexPath)
}
