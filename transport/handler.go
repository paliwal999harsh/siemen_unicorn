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
		http.Error(w, utils.GetAsJsonString(&model.ApiResponse{Msg: "please give amount in natural numbers"}),
			http.StatusBadRequest)
		return
	}
	reqID := h.unicornRequestService.CreateRequest(amount)
	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write(fmt.Append(nil, utils.GetAsJsonString(&model.RequestAcceptedResponse{ReqId: reqID})))
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
		http.Error(w, utils.GetAsJsonString(&model.ApiResponse{Msg: fmt.Sprintf("invalid req id: %s", reqId)}),
			http.StatusNotFound)
		return
	}
	_, _ = w.Write(fmt.Append(nil, utils.GetAsJsonString(req)))
}

func (h *UnicornHandler) GetUnicorn(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	reqId := model.UnicornRequestId(values.Get("id"))
	req, ok := h.unicornRequestService.GetRequest(reqId)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if req.RequestedAmount == req.ReceivedAmount {
		w.WriteHeader(http.StatusNoContent)
		_, _ = w.Write(fmt.Append(nil, utils.GetAsJsonString(&model.ApiResponse{Msg: "Unicorn Request Completed"})))
		return
	}
	items := h.unicornService.GetUnicorn(reqId)
	if len(items) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_, _ = w.Write(fmt.Append(nil, utils.GetAsJsonString(items)))
}
