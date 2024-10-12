package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	"github.com/amankumarsingh77/google_drive_index/internal/drive"
)

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Q         string `json:"q"`
		PageToken string `json:"page_token"`
		PageIndex int    `json:"page_index"`
		Limit     int    `json:"limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keyword := drive.SearchFunction.FormatSearchKeyword(requestData.Q)
	if keyword == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"nextPageToken": nil,
			"curPageIndex":  requestData.PageIndex,
			"data":          nil,
		})
		return
	}

	// gd, err := drive.NewGoogleDriveWithServiceAccount()
	// if err != nil {
	// 	http.Error(w, "Failed to initialize Google Drive client", http.StatusInternalServerError)
	// 	return
	// }
	gd, err := drive.NewGoogleDrive(config.Auth.ClientID, config.Auth.ClientSecret, config.Auth.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to initialize Google Drive client", http.StatusInternalServerError)
		return
	}

	words := strings.Fields(keyword)
	nameSearchStr := fmt.Sprintf("name contains '%s'", strings.Join(words, "' AND name contains '"))
	query := fmt.Sprintf("trashed = false AND mimeType != 'application/vnd.google-apps.shortcut' AND "+
		"mimeType != 'application/vnd.google-apps.document' AND mimeType != 'application/vnd.google-apps.spreadsheet' AND "+
		"mimeType != 'application/vnd.google-apps.form' AND mimeType != 'application/vnd.google-apps.site' AND "+
		"name !='.password' AND (%s)", nameSearchStr)

	params := map[string]string{
		"q":        query,
		"fields":   "nextPageToken, files(id, driveId, parents, name, mimeType, size, modifiedTime)",
		"pageSize": fmt.Sprintf("%d", config.Auth.SearchResultListPageSize),
		"orderBy":  "folder,name,modifiedTime desc",
	}

	if config.Auth.SearchAllDrives {
		params["corpora"] = "allDrives"
		params["includeItemsFromAllDrives"] = "true"
		params["supportsAllDrives"] = "true"
	} else {
		params["corpora"] = "user"
	}

	if requestData.PageToken != "" {
		params["pageToken"] = requestData.PageToken
	}

	searchResult, err := gd.SearchFiles(params)
	if err != nil {
		http.Error(w, "Search failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encrypt file IDs
	// for i, file := range searchResult.Files {
	// 	encryptedID, err := utils.EncryptString(file.Id)
	// 	if err != nil {
	// 		http.Error(w, "Encryption failed", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	searchResult.Files[i].Id = encryptedID
	// }

	response := map[string]interface{}{
		"nextPageToken": searchResult.NextPageToken,
		"curPageIndex":  requestData.PageIndex,
		"data":          searchResult,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
