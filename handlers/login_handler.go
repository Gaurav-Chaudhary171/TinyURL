package handlers

import (
	"TinyURL_Refactored/config"
	"TinyURL_Refactored/model"
	"encoding/json"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
}

type LoginResponse struct {
	Status string `json:"status"`
	User   struct {
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"user"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var user model.Users
	result := config.DB.Select("username, first_name, last_name").Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	response := LoginResponse{
		Status: "success",
	}
	response.User.Username = user.Username
	response.User.FirstName = user.FirstName
	response.User.LastName = user.LastName

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
