package utils

import (
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func HandleError(w http.ResponseWriter, err error, status int, message string) {
	log.Printf("Error: %v", err)
	http.Error(w, message, status)
}

func LogInfo(message string) {
	log.Printf("Info: %s", message)
}
