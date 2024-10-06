package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	grdrive "github.com/amankumarsingh77/google_drive_index/internal/drive"
)

func HandleAPI(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Path      string `json:"path"`
		PageToken string `json:"page_token"`
		PageIndex int    `json:"page_index"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gd, err := grdrive.NewGoogleDrive(config.Auth.ClientID, config.Auth.ClientSecret, config.Auth.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to initialize Google Drive client", http.StatusInternalServerError)
		return
	}

	folderID, err := gd.GetFolderIDFromPath(requestData.Path)
	if err != nil {
		http.Error(w, "Failed to get folder ID: "+err.Error(), http.StatusInternalServerError)
		return
	}

	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	params := map[string]string{
		"q":        query,
		"fields":   "nextPageToken, files(id, name, mimeType, size, modifiedTime)",
		"pageSize": fmt.Sprintf("%d", config.Auth.FilesListPageSize),
		"orderBy":  "folder,name,modifiedTime desc",
	}

	if requestData.PageToken != "" {
		params["pageToken"] = requestData.PageToken
	}

	fileList, err := gd.SearchFiles(params)
	if err != nil {
		http.Error(w, "Failed to list files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"nextPageToken": fileList.NextPageToken,
		"curPageIndex":  requestData.PageIndex,
		"data":          fileList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
