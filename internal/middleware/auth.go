package middleware

import (
	"context"
	"net/http"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !config.Auth.EnableLogin {
			next.ServeHTTP(w, r)
			return
		}

		cookie, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenInfo, err := validateGCPToken(cookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tokenInfo)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func validateGCPToken(token string) (*oauth2.Tokeninfo, error) {
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithAPIKey(token))
	if err != nil {
		return nil, err
	}
	tokenInfo, err := oauth2Service.Tokeninfo().Do()
	if err != nil {
		return nil, err
	}

	return tokenInfo, nil
}
