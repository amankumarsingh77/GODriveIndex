package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amankumarsingh77/google_drive_index/internal/middleware"
	"github.com/amankumarsingh77/google_drive_index/internal/utils"
)

func HandleAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code parameter", http.StatusBadRequest)
		return
	}

	if middleware.OAuthConfig == nil {
		middleware.InitOAuthConfig()
	}

	token, err := middleware.OAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	if err := utils.SaveToken("token.json", token); err != nil {
		log.Printf("Failed to save token: %v", err)
		http.Error(w, "Failed to save authentication token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "<h1>Authentication successful</h1><p>You can close this window now.</p>")
}
