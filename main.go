package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo/router"

	_ "github.com/lib/pq"
)

func main() {
	ginRouter := router.SetupRouter()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: ginRouter,
	}
	graceful := make(chan os.Signal)
	signal.Notify(graceful, syscall.SIGINT)
	signal.Notify(graceful, syscall.SIGTERM)
	go func() {
		<-graceful
		log.Println("Shutting down server...")
		ctx, cancelFunc := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancelFunc()
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatalf("Could not do graceful shutdown: %v\n", err)
		}
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not do graceful shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped...")
}
