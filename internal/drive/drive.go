package drive

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type GoogleDrive struct {
	service *drive.Service
	config  *oauth2.Config
}

func GetRandomServiceAccount() (*config.ServiceAccount, error) {
	if len(config.SAS) == 0 {
		return nil, fmt.Errorf("no service accounts configured")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomIndex := r.Intn(len(config.SAS))
	selectedAccount := config.SAS[randomIndex]

	log.Printf("Selected service account: %s (index: %d)", selectedAccount.ClientEmail, randomIndex)

	return &selectedAccount, nil
}

func NewGoogleDrive(clientID, clientSecret, refreshToken string) (*GoogleDrive, error) {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveScope},
	}

	// token := &oauth2.Token{
	// 	RefreshToken: refreshToken,
	// }
	tokFile := "token.json"
	token, _ := tokenFromFile(tokFile)

	client := config.Client(context.Background(), token)

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &GoogleDrive{
		service: srv,
		config:  config,
	}, nil
}

func NewGoogleDriveWithServiceAccount() (*GoogleDrive, error) {
	serviceAccount, err := GetRandomServiceAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get random service account: %v", err)
	}

	serviceAccountJSON, err := json.Marshal(serviceAccount)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service account: %v", err)
	}

	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithCredentialsJSON(serviceAccountJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create Google Drive service: %v", err)
	}

	return &GoogleDrive{
		service: driveService,
	}, nil
}

func (gd *GoogleDrive) ListFiles(parentID string, pageToken string) (*drive.FileList, error) {
	query := fmt.Sprintf("'%s' in parents and trashed = false", parentID)
	call := gd.service.Files.List().Q(query).
		Fields("nextPageToken, files(id, name, mimeType, size, modifiedTime)")

	if pageToken != "" {
		call = call.PageToken(pageToken)
	}

	return call.Do()
}

func (gd *GoogleDrive) GetFile(fileID string) (*drive.File, error) {
	return gd.service.Files.Get(fileID).
		Fields("id, name, mimeType, size, modifiedTime, webContentLink").
		Do()
}

func (gd *GoogleDrive) SearchFiles(params map[string]string) (*drive.FileList, error) {
	call := gd.service.Files.List()

	if query, ok := params["q"]; ok && query != "" {
		call = call.Q(query)
	}

	if fields, ok := params["fields"]; ok {
		call = call.Fields(googleapi.Field(fields))
	} else {
		call = call.Fields("nextPageToken, files(id, name, mimeType, size, modifiedTime)")
	}

	if pageSize, ok := params["pageSize"]; ok {
		if size, err := strconv.ParseInt(pageSize, 10, 64); err == nil {
			call = call.PageSize(size)
		}
	}

	if pageToken, ok := params["pageToken"]; ok {
		call = call.PageToken(pageToken)
	}

	if orderBy, ok := params["orderBy"]; ok {
		call = call.OrderBy(orderBy)
	}

	call = call.IncludeItemsFromAllDrives(true)
	call = call.SupportsAllDrives(true)

	if driveID, ok := params["driveId"]; ok {
		call = call.DriveId(driveID)
		call = call.Corpora("drive")
	} else {
		call = call.Corpora("allDrives")
	}

	return call.Do()
}

func (gd *GoogleDrive) DownloadFile(fileID string, rangeHeader string) (*drive.File, io.ReadCloser, int64, error) {
	file, err := gd.service.Files.Get(fileID).SupportsAllDrives(true).Fields("id, name, mimeType, size, modifiedTime").Do()
	if err != nil {
		return nil, nil, 0, err
	}

	req := gd.service.Files.Get(fileID).SupportsAllDrives(true)
	if rangeHeader != "" {
		req.Header().Set("Range", rangeHeader)
	}

	resp, err := req.Download()
	if err != nil {
		return nil, nil, 0, err
	}

	return file, resp.Body, resp.ContentLength, nil
}

func (gd *GoogleDrive) GetPasswordForPath(path string) (string, error) {
	parentID, err := gd.GetFolderIDFromPath(path)
	if err != nil {
		return "", err
	}

	query := fmt.Sprintf("'%s' in parents and name = '.password' and trashed = false", parentID)
	fileList, err := gd.service.Files.List().Q(query).Fields("files(id)").Do()
	if err != nil {
		return "", err
	}

	if len(fileList.Files) == 0 {
		return "", nil
	}

	passwordFileID := fileList.Files[0].Id
	content, err := gd.DownloadTextFile(passwordFileID)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(content), nil
}

func (gd *GoogleDrive) GetFolderIDFromPath(path string) (string, error) {
	if !isValidPath(path) {
		return "", errors.New("invalid path")
	}

	parts := splitPath(path)
	currentID := "root"

	for _, part := range parts {
		if part == "" {
			continue
		}

		query := formatQuery(currentID, part, DriveFixedTerms.FolderMimeType)
		fileList, err := gd.service.Files.List().Q(query).Fields("files(id)").Do()
		if err != nil {
			return "", err
		}

		if len(fileList.Files) == 0 {
			return "", fmt.Errorf("folder not found: %s", part)
		}

		currentID = fileList.Files[0].Id
	}

	return currentID, nil
}

func (gd *GoogleDrive) DownloadTextFile(fileID string) (string, error) {
	resp, err := gd.service.Files.Get(fileID).Download()
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
