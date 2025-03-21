package transport

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicorn/model"
	"unicorn/service"
	"unicorn/utils"
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	reqId := model.UnicornRequestId(pathParts[5])
	req, ok := h.unicornRequestService.GetRequest(reqId)
	if !ok {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: fmt.Sprintf("invalid req id: %s", reqId)},
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
	reqId := model.UnicornRequestId(values.Get("id"))
	req, ok := h.unicornRequestService.GetRequest(reqId)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if req.RequestedAmount == req.ReceivedAmount {
		utils.WriteJsonResponseWithStatus(w,
			&model.ApiResponse{Msg: "Unicorn Request Completed"},
			http.StatusNoContent)
		return
	}
	items := h.unicornService.GetUnicorn(reqId)
	if len(items) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	utils.WriteJsonResponse(w, items)
}
