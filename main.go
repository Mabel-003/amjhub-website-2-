package main

import (
	"embed"
	iofs "io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"amjhub/backend"
)

//go:embed static/*
var embeddedAssets embed.FS

func main() {
	// Load local .env for development if present (does not override existing env vars)
	if err := loadDotEnv(".env"); err != nil {
		// non-fatal: proceed if there's no .env
		log.Printf("note: .env not loaded: %v", err)
	}

	mux := http.NewServeMux()

	// Templates are auto-loaded via backend.init() from embedded assets.

	// Serve static assets (CSS, JS, images) from embedded assets
	sub, err := iofs.Sub(embeddedAssets, "static")
	if err == nil {
		fileServer := http.FileServer(http.FS(sub))
		mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	} else {
		// fallback to disk if embedding not available
		fs := http.FileServer(http.Dir("./static"))
		mux.Handle("/static/", http.StripPrefix("/static/", fs))
	}

	// Page routes
	mux.HandleFunc("/", backend.HomeHandler)
	mux.HandleFunc("/work", backend.PortfolioHandler)
	mux.HandleFunc("/portfolio", backend.PortfolioHandler) // alias

	// Form submission route
	mux.HandleFunc("/contact", backend.ContactHandler)

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	addr := ":" + port
	log.Printf("🎬 AMJ HUB server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// loadDotEnv is a tiny helper that loads KEY=VALUE lines from a file
// into the process environment if the key is not already present. It
// intentionally does not support advanced features — it's only for local
// development convenience so users can keep a .env file without an
// external dependency.
func loadDotEnv(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Strip surrounding quotes if present
		if len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}
		if os.Getenv(key) == "" {
			os.Setenv(key, val)
		}
	}
	return nil
}
