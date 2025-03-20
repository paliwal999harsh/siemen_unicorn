package middleware

import (
	"log"
	"net/http"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("| %15s | [%7s] | %-30s | %-80s | from %s", r.Host, r.Method, r.URL.Path, r.URL.RawQuery, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
