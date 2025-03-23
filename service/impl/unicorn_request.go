package impl

import (
	"fmt"
	"time"
	"unicorn/model"
	"unicorn/service"
	"unicorn/storage"
)

type unicornRequestService struct {
	unicornRequestTracker storage.RequestTracker
}

func NewUnicornRequestService(unicornRequestTracker storage.RequestTracker) service.UnicornRequestService {
	return &unicornRequestService{unicornRequestTracker: unicornRequestTracker}
}

func (u *unicornRequestService) GetRequest(reqId model.UnicornRequestId) (*model.UnicornRequest, bool) {
	req, ok := u.unicornRequestTracker.GetRequest(reqId)
	if !ok {
		return &model.UnicornRequest{}, false
	}
	return req, true
}

func (u *unicornRequestService) CreateRequest(amount int) model.UnicornRequestId {
	reqId := model.UnicornRequestId(fmt.Sprintf("REQ-%d", time.Now().Unix()))
	u.unicornRequestTracker.CreateRequest(reqId, amount)
	return reqId
}
