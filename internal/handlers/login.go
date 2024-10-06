package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/amankumarsingh77/google_drive_index/internal/auth"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var serviceAccountData map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&serviceAccountData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	serviceAccountJSON, err := json.Marshal(serviceAccountData)
	if err != nil {
		http.Error(w, "Failed to process service account data", http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateGCPToken(string(serviceAccountJSON))
	if err != nil {
		http.Error(w, "Failed to generate GCP token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Expires:  time.Now().Add(time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
