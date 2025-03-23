package impl

import (
	"unicorn/internal/service"
	"unicorn/internal/storage"
	"unicorn/pkg/model"
)

type unicornService struct {
	unicornStore          storage.UnicornStore
	unicornRequestTracker storage.RequestTracker
}

func NewUnicornService(unicornStore storage.UnicornStore,
	unicornRequestTracker storage.RequestTracker) service.UnicornService {
	return &unicornService{
		unicornStore:          unicornStore,
		unicornRequestTracker: unicornRequestTracker}
}

func (s *unicornService) GetUnicorn(reqId model.UnicornRequestId) []model.Unicorn {
	req, ok := s.unicornRequestTracker.GetRequest(reqId)
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
		s.unicornRequestTracker.UpdateRequest(reqId, req)
		return s.unicornStore.GetUnicorns(int(take))
	}
	return nil
}
