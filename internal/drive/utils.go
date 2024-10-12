package drive

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

func formatQuery(parentID, name, mimeType string) string {
	return fmt.Sprintf("'%s' in parents and name = '%s' and mimeType = '%s' and trashed = false", parentID, name, mimeType)
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}

func isValidPath(path string) bool {
	return path != "" && !strings.ContainsAny(path, "<>:\"|?*")
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func VerifyDownloadToken(token, fileID string) bool {
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
