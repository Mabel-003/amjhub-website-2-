package handler

import (
	"net/http"
	"os"

	"amjhub/backend"
)

func init() {
	if err := backend.LoadTemplates(os.DirFS(".")); err != nil {
		panic(err)
	}
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data := backend.PageData{ActiveNav: "portfolio"}
	if err := backend.ExecutePortfolioTemplate(w, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
