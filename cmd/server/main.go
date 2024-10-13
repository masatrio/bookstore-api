package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masatrio/bookstore-api/config" // Import the config package
	handler "github.com/masatrio/bookstore-api/internal/delivery/http"
	"github.com/masatrio/bookstore-api/utils"
)

func main() {
	ctx := context.Background()

	// Load the application configuration
	cfg := config.LoadConfig()

	tracer := utils.NewTracer(ctx, cfg.Server.ServiceName)

	router := handler.InitRoutes(tracer)
	ServeHTTP(router, cfg.Server.Port)
}

// ServeHTTP serve HTTP API gracefully
func ServeHTTP(router http.Handler, port int) {
	srv := &http.Server{
		Addr:    toPort(port),
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Starting server on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}

func toPort(port int) string {
	return fmt.Sprintf(":%d", port)
}
