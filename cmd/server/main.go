package main

import (
	"log"
	"net/http"
	"os"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
	"github.com/amankumarsingh77/google_drive_index/internal/handlers"
	"github.com/amankumarsingh77/google_drive_index/internal/middleware"
	"github.com/gorilla/mux"
)

func main() {

	config.LoadConfig()

	r := mux.NewRouter()

	r.HandleFunc("/", middleware.AuthMiddleware(handlers.HandleHome))

	r.HandleFunc("/search", middleware.AuthMiddleware(handlers.HandleSearch))
	r.HandleFunc("/download/{id:.+}", handlers.HandleDownload).Queries("token", "{token}")
	r.HandleFunc("/list-files", middleware.AuthMiddleware(handlers.HandleListDriveContents)).Methods("POST")
	r.HandleFunc("/auth/callback", handlers.HandleAuthCallback).Methods("GET")
	r.HandleFunc("/generate-download-link", middleware.AuthMiddleware(handlers.HandleGenerateDownloadLink)).Methods("POST")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
