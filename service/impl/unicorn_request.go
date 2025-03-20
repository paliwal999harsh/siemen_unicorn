package impl

import (
	"fmt"
	"time"
	"unicorn/model"
	"unicorn/service"
	"unicorn/storage"
)

type unicornRequestService struct {
	unicornRequestStorage *storage.InMemoryRequestTracker
}

func NewUnicornRequestService(unicornRequestStorage *storage.InMemoryRequestTracker) service.UnicornRequestService {
	return &unicornRequestService{unicornRequestStorage: unicornRequestStorage}
}

func (u *unicornRequestService) GetRequest(reqId model.UnicornRequestId) (model.UnicornRequest, bool) {
	req, ok := u.unicornRequestStorage.GetRequest(reqId)
	if !ok {
		return model.UnicornRequest{}, false
	}
	return req, true
}

func (u *unicornRequestService) CreateRequest(amount int) model.UnicornRequestId {
	reqId := model.UnicornRequestId(fmt.Sprintf("REQ-%d", time.Now().Unix()))
	u.unicornRequestStorage.CreateRequest(reqId, amount)
	go func(reqId model.UnicornRequestId) {
		req, _ := u.unicornRequestStorage.GetRequest(reqId)
		req.Status = model.UnicornRequestInProgress
		u.unicornRequestStorage.UpdateRequest(reqId, req)
	}(reqId)
	return reqId
}
