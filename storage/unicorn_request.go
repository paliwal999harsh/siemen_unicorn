package storage

import (
	"sync"
	"unicorn/model"
)

type RequestTracker interface {
	CreateRequest(model.UnicornRequestId, int)
	GetRequest(model.UnicornRequestId) (*model.UnicornRequest, bool)
	UpdateRequest(model.UnicornRequestId, *model.UnicornRequest) bool
	GetNextRequest() (model.UnicornRequestId, *model.UnicornRequest, bool)
	RequeueRequest(model.UnicornRequestId, *model.UnicornRequest)
}
type InMemoryRequestTracker struct {
	requestQueue []model.UnicornRequestId
	requests     map[model.UnicornRequestId]*model.UnicornRequest
	sync.Mutex
}

func (rt *InMemoryRequestTracker) GetNextRequest() (model.UnicornRequestId, *model.UnicornRequest, bool) {
	rt.Lock()
	defer rt.Unlock()

	if len(rt.requestQueue) == 0 {
		return model.UnicornRequestId(rune(0)), &model.UnicornRequest{}, false
	}

	reqId := rt.requestQueue[0]
	rt.requestQueue = rt.requestQueue[1:]
	req := rt.requests[reqId]
	return reqId, req, true
}

func (rt *InMemoryRequestTracker) RequeueRequest(id model.UnicornRequestId, req *model.UnicornRequest) {
	rt.Lock()
	defer rt.Unlock()
	rt.requestQueue = append(rt.requestQueue, id)
	rt.requests[id] = req
}

func NewInMemoryRequestTracker() RequestTracker {
	return &InMemoryRequestTracker{
		requests: make(map[model.UnicornRequestId]*model.UnicornRequest),
	}
}

func (rt *InMemoryRequestTracker) CreateRequest(id model.UnicornRequestId, amount int) {
	rt.Lock()
	defer rt.Unlock()
	req := model.UnicornRequest{
		Status:          model.UnicornRequestQueued,
		RequestedAmount: amount,
	}
	rt.requestQueue = append(rt.requestQueue, id)
	rt.requests[id] = &req
}

func (rt *InMemoryRequestTracker) GetRequest(id model.UnicornRequestId) (*model.UnicornRequest, bool) {
	rt.Lock()
	defer rt.Unlock()
	if val, exists := rt.requests[id]; exists {
		return val, true
	}
	return &model.UnicornRequest{}, false
}

func (rt *InMemoryRequestTracker) UpdateRequest(id model.UnicornRequestId, req *model.UnicornRequest) bool {
	rt.Lock()
	defer rt.Unlock()
	if _, exists := rt.requests[id]; !exists {
		return false
	}
	rt.requests[id] = req
	return true
}
