package middleware

import (
	"log"
	"log/slog"
	"net/http"
)

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
		log.Printf(
			"%s %s %s",
			r.Method,
			slog.String("path", r.RequestURI),
			slog.String("remote_addr", r.RemoteAddr),
		)
	}
}
