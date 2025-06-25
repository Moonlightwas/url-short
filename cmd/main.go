package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"url_short/internal/config"
	"url_short/internal/database"
	"url_short/internal/server"
	"url_short/internal/server/handlers"
	"url_short/internal/server/middleware"
)

func main() {
	cfg := config.MustLoad()

	params := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSL)

	db, err := database.Init(params)
	if err != nil {
		slog.Error("failed to initialize database", err)
		return
	}

	log.Printf("starting server on address %s\n", cfg.Server.Address)

	router := server.NewRouter()
	router.Use(middleware.Logger)

	handlers := handlers.NewHandlers(db)

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.Index(w, r)
		case http.MethodPost:
			handlers.UrlHandler(w, r)
		}
	})

	router.HandleFunc("/{alias}", handlers.Redirect)

	//static files
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("Server failed", "error", err)
	}
}
