package handlers

import (
	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
	"encoding/json"
	"net/http"
)

type ShortenRequest struct {
	URL      string `json:"url"`
	Username string `json:"username"`
}

type ShortenResponse struct {
	Status      string `json:"status"`
	ShortenURL  string `json:"shortenurl"`
	OriginalURL string `json:"originalurl"`
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" || req.Username == "" {
		http.Error(w, "URL and username are required", http.StatusBadRequest)
		return
	}

	shortURL := ShortenURL(req.URL)
	shortURL = "https://" + shortURL

	url := &model.GeneratedUrl{
		OriginalUrl: req.URL,
		Username:    req.Username,
		TinyUrl:     shortURL,
		IsActive:    true,
	}

	if err := config.DB.Create(url).Error; err != nil {
		http.Error(w, "Error Inserting Data into DB, may be URL already exists", http.StatusInternalServerError)
		return
	}

	response := ShortenResponse{
		Status:      "success",
		ShortenURL:  shortURL,
		OriginalURL: req.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
