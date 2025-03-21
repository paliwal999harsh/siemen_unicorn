package storage

import (
	"sync"
	"unicorn/model"
)

type RequestTracker interface {
	CreateRequest(model.UnicornRequestId, int)
	GetRequest(model.UnicornRequestId) (model.UnicornRequest, bool)
	UpdateRequest(model.UnicornRequestId, model.UnicornRequest) bool
}
type InMemoryRequestTracker struct {
	requests sync.Map
}

func NewInMemoryRequestTracker() RequestTracker {
	return &InMemoryRequestTracker{}
}

func (rt *InMemoryRequestTracker) CreateRequest(id model.UnicornRequestId, amount int) {
	req := model.UnicornRequest{
		Status:          model.UnicornRequestQueued,
		RequestedAmount: amount,
	}
	rt.requests.Store(id, req)
}

func (rt *InMemoryRequestTracker) GetRequest(id model.UnicornRequestId) (model.UnicornRequest, bool) {
	if val, exists := rt.requests.Load(id); exists {
		data := val.(model.UnicornRequest)
		return data, true
	}
	return model.UnicornRequest{}, false
}

func (rt *InMemoryRequestTracker) UpdateRequest(id model.UnicornRequestId, req model.UnicornRequest) bool {
	_, ok := rt.GetRequest(id)
	if !ok {
		return false
	}
	rt.requests.Store(id, req)
	return true
}
