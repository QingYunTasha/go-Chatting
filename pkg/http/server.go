package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Server defines the HTTP server.
type Server struct {
	httpServer *http.Server
}

// NewServer creates a new instance of the HTTP server.
func NewServer(addr string, handler http.Handler) *Server {
	s := &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}

	return s
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("error starting HTTP server: %v", err)
		}
	}()

	fmt.Printf("HTTP server listening on %s\n", s.httpServer.Addr)

	// Listen for interrupt signals to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	fmt.Printf("Received signal: %s, shutting down server...\n", sig)

	// Gracefully shut down the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("error shutting down HTTP server: %v", err)
	}

	fmt.Println("HTTP server stopped.")
	return nil
}
