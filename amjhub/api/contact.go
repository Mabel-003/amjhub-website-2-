package handler

import (
	"net/http"

	"amjhub/backend"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	backend.ContactHandler(w, r)
}
