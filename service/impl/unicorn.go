package impl

import (
	"unicorn/model"
	"unicorn/service"
	"unicorn/storage"
)

type unicornService struct {
	unicornStore          storage.UnicornStore
	unicornRequestStorage storage.RequestTracker
}

func NewUnicornService(unicornStore storage.UnicornStore,
	unicornRequestStorage storage.RequestTracker) service.UnicornService {
	return &unicornService{
		unicornStore:          unicornStore,
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
	if req.Status == model.UnicornRequestInProgress {
		take := req.AvailableAmount.Load()
		req.AvailableAmount.Store(0)
		req.ReceivedAmount.Add(take)
		if int(req.ReceivedAmount.Load()) == req.RequestedAmount {
			req.Status = model.UnicornRequestCompleted
		}
		s.unicornRequestStorage.UpdateRequest(reqId, req)
		return s.unicornStore.GetUnicorns(int(take))
	}
	return nil
}
