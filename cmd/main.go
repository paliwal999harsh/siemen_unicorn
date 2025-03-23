package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"unicorn/internal/common/middleware"
	"unicorn/internal/factory"
	"unicorn/internal/service/impl"
	"unicorn/internal/storage"
	"unicorn/internal/transport"
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
