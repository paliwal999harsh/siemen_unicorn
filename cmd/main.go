package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unicorn/middleware"
	"unicorn/service/impl"
	"unicorn/storage"
	"unicorn/transport"
)

func main() {
	unicornProducer := storage.NewRandomUnicornProducer()
	unicornStorage := storage.NewInMemoryUnicornStorage()
	unicornRequestStorage := storage.NewInMemoryRequestTracker()

	go unicornSupplier(unicornStorage, unicornProducer)

	unicornService := impl.NewUnicornService(unicornProducer, unicornStorage, &unicornRequestStorage)
	unicornRequestService := impl.NewUnicornRequestService(&unicornRequestStorage)

	unicornHandler := transport.NewUnicornHandler(unicornService, unicornRequestService)

	mux := http.NewServeMux()
	transport.RegisterHealthCheckRoute(mux)
	transport.RegisterUnicornRoutes(mux, unicornHandler)
	wrappedMux := middleware.LoggerMiddleware(middleware.JsonMiddleware(mux))

	setupServer(wrappedMux)
	log.Println("Server started successfully...")
}

func unicornSupplier(storage *storage.InMemoryUnicornStorage, producer storage.UnicornProducer) {
	for {
		storage.SaveUnicorn(producer.CreateUnicorn())
		time.Sleep(5 * time.Second)
	}
}

func setupServer(mux http.Handler) {
	server := &http.Server{Addr: ":8888", Handler: mux}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-quit
		fmt.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			fmt.Println("Server forced to shutdown:", err)
		}
		fmt.Println("Server exited gracefully")
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Server failed to start:", err)
	}
}
