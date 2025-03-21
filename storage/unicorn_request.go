package storage

import (
	"sync"
	"unicorn/model"
)

type RequestTracker interface {
	CreateRequest(model.UnicornRequestId, int)
	GetRequest(model.UnicornRequestId) (model.UnicornRequest, bool)
	UpdateRequest(model.UnicornRequestId, model.UnicornRequest) bool
	GetNextRequest() (model.UnicornRequestId, model.UnicornRequest, bool)
	RequeueRequest(model.UnicornRequestId, model.UnicornRequest)
}
type InMemoryRequestTracker struct {
	requests     []model.UnicornRequestId
	requestsData map[model.UnicornRequestId]model.UnicornRequest
	sync.Mutex
}

func (rt *InMemoryRequestTracker) GetNextRequest() (model.UnicornRequestId, model.UnicornRequest, bool) {
	rt.Lock()
	defer rt.Unlock()

	if len(rt.requests) == 0 {
		return model.UnicornRequestId(rune(0)), model.UnicornRequest{}, false
	}

	reqId := rt.requests[0]
	rt.requests = rt.requests[1:]
	return reqId, rt.requestsData[reqId], true
}

func (rt *InMemoryRequestTracker) RequeueRequest(id model.UnicornRequestId, req model.UnicornRequest) {
	rt.Lock()
	defer rt.Unlock()
	rt.requests = append(rt.requests, id)
	rt.requestsData[id] = req
}

func NewInMemoryRequestTracker() RequestTracker {
	return &InMemoryRequestTracker{
		requestsData: make(map[model.UnicornRequestId]model.UnicornRequest),
	}
}

func (rt *InMemoryRequestTracker) CreateRequest(id model.UnicornRequestId, amount int) {
	rt.Lock()
	defer rt.Unlock()
	req := model.UnicornRequest{
		Status:          model.UnicornRequestQueued,
		RequestedAmount: amount,
	}
	rt.requests = append(rt.requests, id)
	rt.requestsData[id] = req
}

func (rt *InMemoryRequestTracker) GetRequest(id model.UnicornRequestId) (model.UnicornRequest, bool) {
	rt.Lock()
	defer rt.Unlock()
	if val, exists := rt.requestsData[id]; exists {
		return val, true
	}
	return model.UnicornRequest{}, false
}

func (rt *InMemoryRequestTracker) UpdateRequest(id model.UnicornRequestId, req model.UnicornRequest) bool {
	rt.Lock()
	defer rt.Unlock()
	if _, exists := rt.requestsData[id]; !exists {
		return false
	}
	rt.requestsData[id] = req
	return true
}
