package handlers

import (
	"html/template"
	"net/http"

	"github.com/amankumarsingh77/google_drive_index/internal/config"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		SiteName string
	}{
		SiteName: config.Auth.SiteName,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
