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

func NewUnicornService(unicornStorage storage.UnicornStore,
	unicornRequestStorage storage.RequestTracker) service.UnicornService {
	return &unicornService{
		unicornStore:          unicornStorage,
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
		take := req.AvailableAmount
		req.AvailableAmount = 0
		req.ReceivedAmount += take
		if req.ReceivedAmount == req.RequestedAmount {
			req.Status = model.UnicornRequestCompleted
		}
		s.unicornRequestStorage.UpdateRequest(reqId, req)
		return s.unicornStore.GetUnicorns(take)
	}
	return nil
}
