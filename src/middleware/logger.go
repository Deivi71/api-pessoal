package middleware

import (
	"log"
	"net/http"
	"time"
)

// Logger registra informações sobre a requisição HTTP
func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	}
}
