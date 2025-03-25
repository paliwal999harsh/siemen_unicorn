package main

import (
	"context"
	"log"
	"sync"
	"time"
	"unicorn/internal/factory"
	"unicorn/internal/storage"
	"unicorn/pkg/model"
)

func unicornSupplier(ctx context.Context, wg *sync.WaitGroup, store storage.UnicornStore, factory factory.UnicornFactory) {
	defer wg.Done()

	ticker := time.NewTicker(UnicornProductionInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Unicorn supplier stopped...")
			return
		case <-ticker.C:
			if !store.IsAtCapacity() {
				unicorn := factory.CreateUnicorn()
				log.Println("Unicorn Created...", unicorn)
				store.SaveUnicorn(unicorn)
			}
		}
	}
}

func unicornRequestProcessor(ctx context.Context, wg *sync.WaitGroup, unicornStore storage.UnicornStore, unicornRequestTracker storage.RequestTracker) {
	defer wg.Done()

	ticker := time.NewTicker(RequestProcessingInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Unicorn request processor stopped...")
			return
		case <-ticker.C:
			reqId, req, ok := unicornRequestTracker.GetNextRequest()
			if !ok {
				continue
			}
			if req.Status != model.UnicornRequestCompleted {
				if req.AvailableAmount.Load() >= BatchProduction {
					continue
				}
				unicornsAvailable := unicornStore.Capacity()
				if unicornsAvailable > 0 {
					take := min(BatchProduction, req.RequestedAmount-int(req.ReceivedAmount.Load()),
						BatchProduction-int(req.AvailableAmount.Load()), unicornsAvailable)
					unicornStore.DecreaseCapacity(take)
					req.AvailableAmount.Add(int32(take))
					req.Status = model.UnicornRequestInProgress
					unicornRequestTracker.UpdateRequest(reqId, req)
					log.Printf("Fulfilling request: %s, Have: %d, Given %d/%d\n", reqId, req.AvailableAmount.Load(), req.ReceivedAmount.Load(), req.RequestedAmount)
				}
				if int(req.ReceivedAmount.Load()+req.AvailableAmount.Load()) < req.RequestedAmount {
					unicornRequestTracker.RequeueRequest(reqId, req)
				}
			}
		}
	}
}
