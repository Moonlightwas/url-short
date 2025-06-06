package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"
	"url_short/internal/database"
)

func Index(w http.ResponseWriter, r *http.Request, db *database.DB) {
	tmpl := template.Must(template.ParseFiles(
		filepath.Join("web", "templates", "index.html"),
	))

	data := struct {
		url string
	}{
		url: "short url",
	}

	err := tmpl.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		slog.Error("failed to render template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func UrlHandler(w http.ResponseWriter, r *http.Request, db *database.DB) {
	var request struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Wrong JSON", http.StatusBadRequest)
		return
	}

	if request.URL == "" {
		log.Println("Url can't be empty", http.StatusBadRequest)
		http.Error(w, "Url can't be empty", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(request.URL, "https://") && !strings.HasPrefix(request.URL, "http://") {
		log.Println("Field must contain a url")
		http.Error(w, "Field must contain a url", http.StatusBadRequest)
		return
	}

	processedURL, err := db.SaveURL(request.URL)
	if err != nil {
		log.Println("Error in database", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"processedURL": "https://" + r.Host + "/" + processedURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Redirect(w http.ResponseWriter, r *http.Request, db *database.DB) {
	url, err := db.GetURL(strings.Trim(r.URL.Path, "/"))
	if err == database.ErrURLNotFound {
		tmpl := template.Must(template.ParseFiles(
			filepath.Join("web", "templates", "404.html"),
		))
		data := struct {
			url string
		}{
			url: "short url",
		}

		err := tmpl.ExecuteTemplate(w, "404.html", data)
		if err != nil {
			slog.Error("failed to render template", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
