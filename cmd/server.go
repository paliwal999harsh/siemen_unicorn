package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

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
