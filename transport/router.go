package transport

import (
	"encoding/json"
	"net/http"
	"unicorn/model"
)

func RegisterHealthCheckRoute(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(model.ApiResponse{Msg: "OK"})
	})
}

func RegisterUnicornRoutes(mux *http.ServeMux, unicornHandler UnicornHandler) {
	mux.HandleFunc("/api/v1/unicorn", unicornHandler.GetUnicorn)
	mux.HandleFunc("/api/v1/unicorn/request", unicornHandler.RequestUnicorn)
	mux.HandleFunc("/api/v1/unicorn/request/", unicornHandler.CheckRequestStatus)
}
