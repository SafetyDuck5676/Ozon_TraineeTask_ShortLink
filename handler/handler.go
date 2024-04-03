// handler.go

package handler

import (
	"Ozon/db"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	storage db.Storage
}

func NewHandler(storage db.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	shortLink, err := h.storage.Shorten(requestData.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	responseData := struct {
		ShortLink string `json:"shortLink"`
	}{ShortLink: shortLink}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}

func (h *Handler) ExpandURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortLink := vars["shortLink"]

	originalURL, err := h.storage.Expand(shortLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
