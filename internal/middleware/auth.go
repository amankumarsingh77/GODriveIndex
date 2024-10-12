package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	OAuthConfig *oauth2.Config
	ErrNoToken  = errors.New("access token is missing or empty")
)

func InitOAuthConfig() *oauth2.Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	OAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/drive.appdata",
			"https://www.googleapis.com/auth/drive.file",
			"https://www.googleapis.com/auth/drive.metadata",
			"https://www.googleapis.com/auth/drive.metadata.readonly",
			"https://www.googleapis.com/auth/drive.photos.readonly",
			"https://www.googleapis.com/auth/drive.readonly",
			"https://www.googleapis.com/auth/drive.scripts",
		},
		Endpoint: google.Endpoint,
	}
	return OAuthConfig
}

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if OAuthConfig == nil {
			OAuthConfig = InitOAuthConfig()
		}

		token, err := validateTokenFile("token.json")
		if err != nil {
			handleAuthError(w, err)
			return
		}

		tokenInfo, err := validateGCPToken(r.Context(), token.RefreshToken)
		if err != nil {
			handleAuthError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, tokenInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func handleAuthError(w http.ResponseWriter, err error) {
	log.Printf("Auth error: %v", err)
	authURL := OAuthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	response := map[string]string{
		"message": "Token is missing or invalid. Please authorize again.",
		"authUrl": authURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(response)
}

func validateGCPToken(ctx context.Context, refreshToken string) (*oauth2.Token, error) {
	if OAuthConfig == nil {
		OAuthConfig = InitOAuthConfig()
	}

	log.Println(refreshToken)

	token, err := OAuthConfig.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken}).Token()
	if err != nil {
		return nil, fmt.Errorf("error refreshing token: %w", err)
	}

	return token, nil
}

func validateTokenFile(filePath string) (*oauth2.Token, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading token file: %w", err)
	}

	var token oauth2.Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("error parsing token file: %w", err)
	}

	if token.AccessToken == "" {
		return nil, ErrNoToken
	}

	return &token, nil
}
