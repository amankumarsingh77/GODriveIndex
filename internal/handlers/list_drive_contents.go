package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	"github.com/amankumarsingh77/google_drive_index/internal/drive"
)

func HandleListDriveContents(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Path      string `json:"path"`
		PageToken string `json:"page_token"`
		PageIndex int    `json:"page_index"`
		Limit     int    `json:"limit"`
		FolderID  string `json:"folder_id,omitempty"`
		DriveID   string `json:"drive_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.Limit <= 0 {
		requestData.Limit = 10
	}

	if requestData.PageIndex < 1 {
		requestData.PageIndex = 1
	}

	gd, err := drive.NewGoogleDrive(config.Auth.ClientID, config.Auth.ClientSecret, config.Auth.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to initialize Google Drive client: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var folderID string
	if requestData.FolderID != "" {
		folderID = requestData.FolderID
	} else {
		folderID, err = gd.GetFolderIDFromPath(requestData.Path)
		if err != nil {
			http.Error(w, "Failed to get folder ID: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	query := fmt.Sprintf("'%s' in parents and trashed = false AND name !='.password' and mimeType != 'application/vnd.google-apps.shortcut' and mimeType != 'application/vnd.google-apps.document' and mimeType != 'application/vnd.google-apps.spreadsheet' and mimeType != 'application/vnd.google-apps.form' and mimeType != 'application/vnd.google-apps.site'", folderID)

	params := map[string]string{
		"q":                         query,
		"fields":                    "nextPageToken, files(id, driveId, name, mimeType, size, modifiedTime, kind, fileExtension)",
		"pageSize":                  fmt.Sprintf("%d", requestData.Limit),
		"orderBy":                   "folder,name,modifiedTime desc",
		"includeItemsFromAllDrives": "true",
		"supportsAllDrives":         "true",
	}

	if requestData.DriveID != "" {
		params["driveId"] = requestData.DriveID
		params["corpora"] = "drive"
	} else {
		params["corpora"] = "allDrives"
	}

	if requestData.PageToken != "" {
		params["pageToken"] = requestData.PageToken
	} else if requestData.PageIndex > 1 {
		for i := 1; i < requestData.PageIndex; i++ {
			tempList, err := gd.SearchFiles(params)
			if err != nil {
				http.Error(w, "Failed to list files: "+err.Error(), http.StatusInternalServerError)
				return
			}
			if tempList.NextPageToken == "" {
				http.Error(w, "Page index out of range", http.StatusBadRequest)
				return
			}
			params["pageToken"] = tempList.NextPageToken
		}
	}

	fileList, err := gd.SearchFiles(params)
	if err != nil {
		log.Printf("Error searching files: %v", err)
		http.Error(w, "Failed to list files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"nextPageToken": fileList.NextPageToken,
		"curPageIndex":  requestData.PageIndex,
		"data":          fileList,
		"limit":         requestData.Limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
