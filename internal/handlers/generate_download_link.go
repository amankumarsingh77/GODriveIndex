package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	"github.com/amankumarsingh77/google_drive_index/internal/utils"
)

func HandleGenerateDownloadLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		FileID string `json:"file_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token := GenerateDownloadToken(requestData.FileID, time.Duration(config.Auth.FileLinkExpiry)*time.Hour)

	encryptedFileID, err := utils.EncryptString(requestData.FileID)
	if err != nil {
		http.Error(w, "Failed to generate download link", http.StatusInternalServerError)
		return
	}

	downloadURL := fmt.Sprintf("%s/download/%s?token=%s", config.Auth.RedirectDomain, encryptedFileID, token)

	response := map[string]string{
		"download_url": downloadURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
