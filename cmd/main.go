package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicorn/middleware"
	"unicorn/model"
	"unicorn/service/impl"
	"unicorn/storage"
	"unicorn/transport"
)

const (
	MaxStoreCapacity          = 100
	BatchProduction           = 10
	UnicornProductionInterval = 5
	RequestProcessingInterval = 5
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	unicornFactory := storage.NewRandomUnicornProducer()
	unicornStore := storage.NewInMemoryUnicornStore()
	unicornRequestTracker := storage.NewInMemoryRequestTracker()

	go unicornSupplier(ctx, unicornStore, unicornFactory)
	go unicornRequestProcessor(ctx, unicornStore, unicornRequestTracker)

	unicornService := impl.NewUnicornService(unicornStore, unicornRequestTracker)
	unicornRequestService := impl.NewUnicornRequestService(unicornRequestTracker)

	unicornHandler := transport.NewUnicornHandler(unicornService, unicornRequestService)

	mux := http.NewServeMux()
	transport.RegisterHealthCheckRoute(mux)
	transport.RegisterUnicornRoutes(mux, unicornHandler)
	wrappedMux := middleware.LoggerMiddleware(middleware.JsonMiddleware(mux))

	setupServer(ctx, wrappedMux)
}

func unicornSupplier(ctx context.Context, store storage.UnicornStore, factory storage.UnicornFactory) {
	ticker := time.NewTicker(UnicornProductionInterval * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Unicorn supplier stopped...")
		case <-ticker.C:
			if store.AvailableUnicorns() < MaxStoreCapacity {
				store.SaveUnicorn(factory.CreateUnicorn())
			}
		}
	}
}

func unicornRequestProcessor(ctx context.Context, unicornStore storage.UnicornStore, unicornRequestTracker storage.RequestTracker) {
	ticker := time.NewTicker(RequestProcessingInterval * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Unicorn processor stopped...")
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

func setupServer(ctx context.Context, mux http.Handler) {
	server := &http.Server{Addr: ":8888", Handler: mux}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Println("Server forced to shutdown:", err)
		}
		log.Println("Server exited gracefully")
	}()
	log.Println("Starting Server, listening on", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("Server failed to start:", err)
	}
}
