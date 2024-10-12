package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
)

func SaveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to create or open token file: %w", err)
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}
