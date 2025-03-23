package utils

import (
	"fmt"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, obj any) {
	_, _ = w.Write(fmt.Append(nil, GetAsJsonString(obj)))
}

func WriteJsonResponseWithStatus(w http.ResponseWriter, obj any, status int) {
	w.WriteHeader(status)
	_, _ = w.Write(fmt.Append(nil, GetAsJsonString(obj)))
}
