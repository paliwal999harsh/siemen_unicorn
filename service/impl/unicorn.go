package impl

import (
	"unicorn/model"
	"unicorn/service"
	"unicorn/storage"
)

type unicornService struct {
	unicornProducer       storage.UnicornProducer
	unicornStorage        storage.UnicornStorage
	unicornRequestStorage storage.RequestTracker
}

func NewUnicornService(unicornProducer storage.UnicornProducer,
	unicornStorage storage.UnicornStorage,
	unicornRequestStorage storage.RequestTracker) service.UnicornService {
	return &unicornService{
		unicornProducer:       unicornProducer,
		unicornStorage:        unicornStorage,
		unicornRequestStorage: unicornRequestStorage}
}

func (s *unicornService) GetUnicorn(reqId model.UnicornRequestId) []model.Unicorn {
	req, ok := s.unicornRequestStorage.GetRequest(reqId)
	if !ok {
		return nil
	}
	if req.Status == model.UnicornRequestQueued {
		return nil
	}
	var data = make(chan []model.Unicorn)
	go func(reqId model.UnicornRequestId, req model.UnicornRequest, data chan []model.Unicorn) {
		var items []model.Unicorn
		amount := req.RequestedAmount - req.ReceivedAmount
		items = s.unicornStorage.GetUnicorns(amount)
		req.ReceivedAmount += len(items)
		if req.RequestedAmount == req.ReceivedAmount {
			req.Status = model.UnicornRequestCompleted
		}
		s.unicornRequestStorage.UpdateRequest(reqId, req)
		data <- items
	}(reqId, req, data)
	return <-data
}
