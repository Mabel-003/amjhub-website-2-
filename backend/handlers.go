package backend

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type ContactFormData struct {
	FullName  string
	Email     string
	Phone     string
	Message   string
	Service   string
	Timestamp time.Time
}

type ContactResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type PageData struct {
	ActiveNav string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var homeTmpl *template.Template
var portfolioTmpl *template.Template

// LoadTemplates parses templates from the provided file system. Call this
// from `main` after embedding assets (or pass the local FS for disk-based
// development).
func LoadTemplates(fsys fs.FS) error {
	var err error
	homeTmpl, err = template.ParseFS(fsys, "templates/base.html", "templates/home.html")
	if err != nil {
		return err
	}
	portfolioTmpl, err = template.ParseFS(fsys, "templates/base.html", "templates/portfolio.html")
	if err != nil {
		return err
	}
	return nil
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageData{ActiveNav: "home"}
	if err := homeTmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("ERROR: failed to execute home template: %v", err)
	}
}

func PortfolioHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := PageData{ActiveNav: "portfolio"}
	if err := portfolioTmpl.ExecuteTemplate(w, "base", data); err != nil {
		log.Printf("ERROR: failed to execute portfolio template: %v", err)
	}
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(ContactResponse{
			Success: false,
			Error:   "Method not allowed. Use POST.",
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ContactResponse{
			Success: false,
			Error:   "Failed to parse form data.",
		})
		return
	}

	data := ContactFormData{
		FullName:  strings.TrimSpace(r.FormValue("full_name")),
		Email:     strings.TrimSpace(r.FormValue("email")),
		Phone:     strings.TrimSpace(r.FormValue("phone")),
		Message:   strings.TrimSpace(r.FormValue("message")),
		Service:   strings.TrimSpace(r.FormValue("service")),
		Timestamp: time.Now(),
	}

	if data.FullName == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ContactResponse{Success: false, Error: "Full name is required."})
		return
	}
	if len(data.FullName) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ContactResponse{Success: false, Error: "Please enter your full name."})
		return
	}
	if data.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ContactResponse{Success: false, Error: "Email address is required."})
		return
	}
	if !emailRegex.MatchString(data.Email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ContactResponse{Success: false, Error: "Please enter a valid email address."})
		return
	}
	if data.Message == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ContactResponse{Success: false, Error: "Please tell us about your project."})
		return
	}

	log.Println("─────────────────────────────────────────────")
	log.Println("📬  NEW CONTACT FORM SUBMISSION — AMJ HUB")
	log.Println("─────────────────────────────────────────────")
	log.Printf("  Timestamp : %s", data.Timestamp.Format("2006-01-02 15:04:05"))
	log.Printf("  Full Name : %s", data.FullName)
	log.Printf("  Email     : %s", data.Email)
	log.Printf("  Phone     : %s", ifEmpty(data.Phone, "(not provided)"))
	log.Printf("  Service   : %s", ifEmpty(data.Service, "(not specified)"))
	log.Printf("  Message   : %s", data.Message)
	log.Println("─────────────────────────────────────────────")

	emailCfg := LoadEmailConfig()
	if !emailCfg.IsBusinessConfigured() {
		logEmailFallback(data, fmt.Errorf("email config incomplete"))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ContactResponse{Success: false, Error: "Email delivery not configured"})
		return
	}

	if err := SendBusinessNotification(emailCfg, data); err != nil {
		logEmailFallback(data, err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ContactResponse{Success: false, Error: "Failed to deliver notification email"})
		return
	}
	log.Printf("✅  Inquiry emailed to %s", emailCfg.ToAddress)

	// Only send a confirmation email to the user after the business email
	// has successfully been delivered.
	if err := SendConfirmationEmail(emailCfg, data); err != nil {
		log.Printf("⚠️  failed to send confirmation email to user %s: %v", data.Email, err)
	}
	log.Println("─────────────────────────────────────────────")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ContactResponse{
		Success: true,
		Message: "Thank you, " + data.FullName + "! Your inquiry has been received. Our team will reach out to you at " + data.Email + " within 24 hours.",
	})
}

func ifEmpty(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

func ExecuteHomeTemplate(w io.Writer, data PageData) error {
	return homeTmpl.ExecuteTemplate(w, "base", data)
}

func ExecutePortfolioTemplate(w io.Writer, data PageData) error {
	return portfolioTmpl.ExecuteTemplate(w, "base", data)
}

func LoadTemplatesFromDisk() error {
	return LoadTemplates(os.DirFS("."))
}
