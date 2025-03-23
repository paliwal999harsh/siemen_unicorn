package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unicorn/factory"
	"unicorn/middleware"
	"unicorn/model"
	"unicorn/service/impl"
	"unicorn/storage"
	"unicorn/transport"
)

const (
	BatchProduction           = 10
	UnicornProductionInterval = 5 * time.Second
	RequestProcessingInterval = 2 * time.Second
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	unicornFactory := factory.NewRandomUnicornProducer()
	unicornStore := storage.NewInMemoryUnicornStore()
	unicornRequestTracker := storage.NewInMemoryRequestTracker()

	wg.Add(2)
	go unicornSupplier(ctx, &wg, unicornStore, unicornFactory)
	go unicornRequestProcessor(ctx, &wg, unicornStore, unicornRequestTracker)

	unicornService := impl.NewUnicornService(unicornStore, unicornRequestTracker)
	unicornRequestService := impl.NewUnicornRequestService(unicornRequestTracker)

	unicornHandler := transport.NewUnicornHandler(unicornService, unicornRequestService)

	mux := http.NewServeMux()
	transport.RegisterHealthCheckRoute(mux)
	transport.RegisterUnicornRoutes(mux, unicornHandler)
	wrappedMux := middleware.LoggerMiddleware(middleware.JsonMiddleware(mux))

	setupServer(wrappedMux, ctx, cancel)
	wg.Wait()
	log.Println("Application shut down successfully")
}

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
			log.Println("Unicorn processor stopped...")
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

func setupServer(mux http.Handler, ctx context.Context, cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	server := &http.Server{Addr: ":8888", Handler: mux}
	go func() {
		<-quit
		defer cancel()
		defer func(server *http.Server, ctx context.Context) {
			if err := server.Shutdown(ctx); err != nil {
				log.Println("Server Shutdown error:", err)
			}
		}(server, ctx)
		log.Println("Shutting down server...")
	}()
	log.Println("Starting Server, listening on", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Println("Server failed to start:", err)
	}
}
