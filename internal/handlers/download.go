package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

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

	if !drive.VerifyDownloadToken(token, fileID) {
		http.Error(w, "Invalid or expired download token", http.StatusForbidden)
		return
	}

	gd, err := drive.NewGoogleDriveWithServiceAccount()
	if err != nil {
		http.Error(w, "Failed to initialize Google Drive client"+err.Error(), http.StatusInternalServerError)
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		_, err = io.CopyN(w, stream, end-start+1)
		if err != nil {
			if strings.Contains(err.Error(), "write: connection reset by peer") || strings.Contains(err.Error(), "write: broken pipe") {
				log.Printf("Client disconnected: %v", err)
			} else {
				log.Printf("Error streaming file: %v", err)
			}
		}
	}()

	wg.Wait()
}
