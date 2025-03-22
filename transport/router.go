package transport

import (
	"net/http"
	"unicorn/model"
	"unicorn/utils"
)

func RegisterHealthCheckRoute(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		utils.WriteJsonResponseWithStatus(w,
			model.ApiResponse{Msg: "OK"},
			http.StatusOK)
	})
}

func RegisterUnicornRoutes(mux *http.ServeMux, unicornHandler UnicornHandler) {
	mux.HandleFunc("/api/v1/unicorn", unicornHandler.GetUnicorn)
	mux.HandleFunc("/api/v1/unicorn/request", unicornHandler.RequestUnicorn)
	mux.HandleFunc("/api/v1/unicorn/request/", unicornHandler.CheckRequestStatus)
}
