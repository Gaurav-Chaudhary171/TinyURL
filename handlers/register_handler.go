package handlers

import (
	"encoding/json"
	"net/http"

	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
)

type RegisterRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RegisterResponse struct {
	Status   string `json:"status"`
	Username string `json:"username"`
}

// generateUsername creates a username from first and last name
func generateUsername(firstName, lastName string) string {
	// Get first half of first name
	firstHalf := firstName
	if len(firstName) > 1 {
		firstHalf = firstName[:len(firstName)/2]
	}

	// Get last half of last name
	lastHalf := lastName
	if len(lastName) > 1 {
		lastHalf = lastName[len(lastName)/2:]
	}

	return firstHalf + lastHalf
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FirstName == "" || req.LastName == "" {
		http.Error(w, "First name and last name are required", http.StatusBadRequest)
		return
	}

	// Generate username from first and last name
	username := generateUsername(req.FirstName, req.LastName)

	// Check if username already exists
	var existingUser model.Users
	result := config.DB.Where("username = ?", username).First(&existingUser)
	if result.Error == nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Create new user
	user := model.Users{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Username:  username,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		http.Error(w, "Error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		Status:   "success",
		Username: user.Username,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
