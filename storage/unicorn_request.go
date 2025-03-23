package storage

import (
	"log"
	"sync"
	"unicorn/model"
	"unicorn/pkg/collection"
)

type RequestTracker interface {
	CreateRequest(model.UnicornRequestId, int)
	GetRequest(model.UnicornRequestId) (*model.UnicornRequest, bool)
	UpdateRequest(model.UnicornRequestId, *model.UnicornRequest) bool
	GetNextRequest() (model.UnicornRequestId, *model.UnicornRequest, bool)
	RequeueRequest(model.UnicornRequestId, *model.UnicornRequest)
}
type InMemoryRequestTracker struct {
	requestQueue collection.Queue[model.UnicornRequestId]
	requests     collection.Map[model.UnicornRequestId, *model.UnicornRequest] // map[model.UnicornRequestId]*model.UnicornRequest
	sync.Mutex
}

func (rt *InMemoryRequestTracker) GetNextRequest() (model.UnicornRequestId, *model.UnicornRequest, bool) {
	rt.Lock()
	defer rt.Unlock()

	if rt.requestQueue.Empty() {
		return model.UnicornRequestId(rune(0)), nil, false
	}

	reqId, err := rt.requestQueue.Poll()
	if err != nil {
		return model.UnicornRequestId(rune(0)), nil, false
	}
	req, ok := rt.requests.Get(reqId)
	if !ok {
		return model.UnicornRequestId(rune(0)), nil, false
	}
	return reqId, req, true
}

func (rt *InMemoryRequestTracker) RequeueRequest(id model.UnicornRequestId, req *model.UnicornRequest) {
	rt.Lock()
	defer rt.Unlock()
	_, _ = rt.requestQueue.Offer(id)
	rt.requests.Put(id, req)
}

func (rt *InMemoryRequestTracker) CreateRequest(id model.UnicornRequestId, amount int) {
	rt.Lock()
	defer rt.Unlock()
	req := model.UnicornRequest{
		Status:          model.UnicornRequestQueued,
		RequestedAmount: amount,
	}
	ok, _ := rt.requestQueue.Offer(id)
	if !ok {
		log.Println("unable to queue the request")
	}
	rt.requests.Put(id, &req)
}

func (rt *InMemoryRequestTracker) GetRequest(id model.UnicornRequestId) (*model.UnicornRequest, bool) {
	rt.Lock()
	defer rt.Unlock()
	if val, exists := rt.requests.Get(id); exists {
		return val, true
	}
	return nil, false
}

func (rt *InMemoryRequestTracker) UpdateRequest(id model.UnicornRequestId, req *model.UnicornRequest) bool {
	rt.Lock()
	defer rt.Unlock()
	if _, exists := rt.requests.Get(id); !exists {
		return false
	}
	rt.requests.Put(id, req)
	return true
}

func NewInMemoryRequestTracker() RequestTracker {
	return &InMemoryRequestTracker{
		requestQueue: collection.NewSliceQueue[model.UnicornRequestId](),
		requests:     collection.NewNativeMap[model.UnicornRequestId, *model.UnicornRequest](),
	}
}
