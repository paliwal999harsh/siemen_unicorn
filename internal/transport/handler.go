package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicorn/internal/service"
	"unicorn/pkg/model"
	"unicorn/pkg/utils"
)

type UnicornHandler struct {
	unicornService        service.UnicornService
	unicornRequestService service.UnicornRequestService
}

func NewUnicornHandler(unicornService service.UnicornService,
	unicornRequestService service.UnicornRequestService) UnicornHandler {
	return UnicornHandler{
		unicornService:        unicornService,
		unicornRequestService: unicornRequestService}
}

func (h *UnicornHandler) RequestUnicorn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	values := r.URL.Query()
	if !values.Has("amount") {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "query mandatory param 'amount' not found."},
			http.StatusBadRequest)
		return
	}
	amount, err := strconv.Atoi(values.Get("amount"))
	if err != nil {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "please give amount in natural numbers"},
			http.StatusBadRequest)
		return
	}
	if amount <= 0 {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "amount must be positive"},
			http.StatusBadRequest)
		return
	}
	reqID := h.unicornRequestService.CreateRequest(amount)
	utils.WriteJsonResponseWithStatus(w,
		&model.RequestAcceptedResponse{ReqId: reqID},
		http.StatusAccepted)
}

func (h *UnicornHandler) CheckRequestStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 6 || pathParts[5] == "" {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "please provide request 'id' in url path"},
			http.StatusBadRequest)
		return
	}
	reqId := model.UnicornRequestId(pathParts[5])
	req, ok := h.unicornRequestService.GetRequest(reqId)
	if !ok {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: fmt.Sprintf("invalid request id: %s", reqId)},
			http.StatusNotFound)
		return
	}
	utils.WriteJsonResponse(w, req)
}

func (h *UnicornHandler) GetUnicorn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	values := r.URL.Query()
	if !values.Has("id") {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "query mandatory parameter request 'id' not found."},
			http.StatusBadRequest)
		return
	}
	reqId := model.UnicornRequestId(values.Get("id"))
	req, ok := h.unicornRequestService.GetRequest(reqId)
	if !ok {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "invalid request 'id', not found."},
			http.StatusNotFound)
		return
	}
	if req.RequestedAmount == int(req.ReceivedAmount.Load()) {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "Unicorn Request Completed"},
			http.StatusOK)
		return
	}
	items := h.unicornService.GetUnicorn(reqId)
	if len(items) == 0 {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "Unicorn not available."},
			http.StatusOK)
		return
	}
	utils.WriteJsonResponse(w, items)
}
