package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	"github.com/amankumarsingh77/google_drive_index/internal/drive"
	"github.com/amankumarsingh77/google_drive_index/internal/utils"
	"github.com/gorilla/mux"
)

func HandleDownload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encryptedFileID := vars["id"]
	token := vars["token"]

	fileID, err := utils.DecryptString(encryptedFileID)
	if err != nil {
		http.Error(w, "Invalid file ID"+err.Error(), http.StatusBadRequest)
		return
	}

	if !verifyDownloadToken(token, fileID) {
		http.Error(w, "Invalid or expired download token", http.StatusForbidden)
		return
	}

	gd, err := drive.NewGoogleDrive(config.Auth.ClientID, config.Auth.ClientSecret, config.Auth.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to initialize Google Drive client", http.StatusInternalServerError)
		return
	}

	rangeHeader := r.Header.Get("Range")
	file, stream, _, err := gd.DownloadFile(fileID, rangeHeader)
	if err != nil {
		http.Error(w, "Failed to download file"+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stream.Close()

	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Type", file.MimeType)

	if rangeHeader == "" {
		w.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))
		w.WriteHeader(http.StatusOK)
		io.Copy(w, stream)
		return
	}

	rangeParts := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
	start, _ := strconv.ParseInt(rangeParts[0], 10, 64)
	end := file.Size - 1
	if len(rangeParts) > 1 && rangeParts[1] != "" {
		end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
	}
	if start > end || start < 0 || end >= file.Size {
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, file.Size))
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.WriteHeader(http.StatusPartialContent)

	_, err = io.CopyN(w, stream, end-start+1)
	if err != nil {
		log.Printf("Error streaming file: %v", err)
	}
}

func verifyDownloadToken(token, fileID string) bool {

	token = strings.Replace(token, " ", "+", -1)
	firstColonIndex := strings.Index(token, ":")
	if firstColonIndex == -1 {
		log.Println("Token format is invalid: no colon found")
		return false
	}

	tokenFileID := token[:firstColonIndex]
	tokenExpiry := token[firstColonIndex+1:]

	expiryTime, err := time.Parse(time.RFC3339, tokenExpiry)
	if err != nil {
		return false
	}

	isValid := tokenFileID == fileID && time.Now().Before(expiryTime)
	log.Printf("Token is valid: %v", isValid)
	return isValid
}

func GenerateDownloadToken(fileID string, duration time.Duration) string {
	expiry := time.Now().Add(duration).Format(time.RFC3339)
	return fmt.Sprintf("%s:%s", fileID, expiry)
}
