package transport

import (
	"net/http"
	"unicorn/middleware"
)

func RegisterHealthCheckRoute() {
	http.Handle("/api/v1/health", middleware.RequestLogger(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})))
}

func RegisterUnicornRoutes(unicornHandler UnicornHandler) {
	http.Handle("/api/v1/unicorn", middleware.RequestLogger(http.HandlerFunc(unicornHandler.GetUnicorn)))
	// http.HandleFunc("/api/v1/unicorn", unicornHandler.GetUnicorn)
}

func RegisterUnicornRequestRoutes(unicornRequestHandler UnicornRequestHandler) {
	http.Handle("/api/v1/unicorn/request", middleware.RequestLogger(http.HandlerFunc(unicornRequestHandler.RequestUnicorn)))
	http.Handle("/api/v1/unicorn/request/", middleware.RequestLogger(http.HandlerFunc(unicornRequestHandler.CheckRequestStatus)))

	// http.HandleFunc("/api/v1/unicorn/request/", unicornRequestHandler.HandleRequestUnicorn)
}
