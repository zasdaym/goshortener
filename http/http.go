package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 120 * time.Second
	shutdownTimeout = 5 * time.Second
)

// Server is an HTTP server.
type Server struct {
	server *http.Server
}

// NewServer creates a new default HTTP server.
func NewServer(addr string) (srv *Server) {
	srv = &Server{
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  readTimeout,
			WriteTimeout: writeTimeout,
			IdleTimeout:  idleTimeout,
		},
	}
	return srv
}

// Serve requests and shutdown gracefully on interrupt.
func (srv *Server) Serve(ctx context.Context) (err error) {
	ch := make(chan error, 1)
	// start HTTP server on another goroutine
	go func() {
		log.Printf("HTTP server listening on %s", srv.server.Addr)
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			ch <- err
		}
	}()

	// wait for interrupt from context or error from previous goroutine
	select {
	case err = <-ch:
		return fmt.Errorf("failed to start HTTP server: %w", err)
	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		log.Printf("shutting down HTTP server gracefully")
		if err := srv.server.Shutdown(ctxShutdown); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server gracefully: %w", err)
		}
		return nil
	}
}
