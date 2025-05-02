package handlers

import (
	"encoding/json"
	"net/http"

	"TinyURL_Refactored/config"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	// Test database connection
	sqlDB, err := config.DB.DB()
	if err != nil {
		respondWithJSON(w, http.StatusServiceUnavailable, HealthResponse{
			Status:  "error",
			Message: "Database connection error: " + err.Error(),
		})
		return
	}

	// Ping the database
	if err := sqlDB.Ping(); err != nil {
		respondWithJSON(w, http.StatusServiceUnavailable, HealthResponse{
			Status:  "error",
			Message: "Database ping failed: " + err.Error(),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, HealthResponse{
		Status:  "healthy",
		Message: "Service is healthy",
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error","message":"Internal server error"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
