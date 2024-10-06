package auth

import (
	"encoding/json"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type ServiceAccount struct {
	ClientEmail string `json:"client_email"`
	PrivateKey  string `json:"private_key"`
}

func GenerateGCPToken(serviceAccountJSON string) (string, error) {
	var sa ServiceAccount
	err := json.Unmarshal([]byte(serviceAccountJSON), &sa)
	if err != nil {
		return "", err
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(sa.PrivateKey))
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.StandardClaims{
		Issuer:    sa.ClientEmail,
		Audience:  "https://oauth2.googleapis.com/token",
		ExpiresAt: now.Add(time.Hour).Unix(),
		IssuedAt:  now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}
