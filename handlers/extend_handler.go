package handlers

import (
	"encoding/json"
	"net/http"

	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
)

type ExtendRequest struct {
	URL      string `json:"url"`
	Username string `json:"username"`
}

type ExtendResponse struct {
	Status      string `json:"status"`
	OriginalURL string `json:"originalurl"`
	ExtendedURL string `json:"extendedurl"`
}

func ExtendHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ExtendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" || req.Username == "" {
		http.Error(w, "URL and username are required", http.StatusBadRequest)
		return
	}

	var url model.GeneratedUrl
	result := config.DB.Where("tiny_url = ? AND username = ? AND is_active = ?", req.URL, req.Username, true).First(&url)
	if result.Error != nil {
		http.Error(w, "Error fetching Data from DB, May be URL is not found or not active", http.StatusNotFound)
		return
	}

	response := ExtendResponse{
		Status:      "success",
		OriginalURL: url.OriginalUrl,
		ExtendedURL: req.URL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
