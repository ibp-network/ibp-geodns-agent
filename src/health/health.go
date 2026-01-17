package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ibp-network/ibp-geodns-agent/src/logging"
)

// Server provides health check endpoints
type Server struct {
	port     int
	server   *http.Server
	mu       sync.RWMutex
	healthy  bool
	ready    bool
	started  bool
}

// New creates a new health server
func New(port int) *Server {
	return &Server{
		port:    port,
		healthy: true,
		ready:   false,
	}
}

// Start starts the health server
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.started {
		return fmt.Errorf("health server already started")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ready", s.readyHandler)
	mux.HandleFunc("/live", s.liveHandler)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Error("Health server error", "error", err)
		}
	}()

	s.started = true
	s.ready = true
	logging.Info("Health server started", "port", s.port)

	return nil
}

// Stop stops the health server
func (s *Server) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.started {
		return nil
	}

	if s.server != nil {
		if err := s.server.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown health server: %w", err)
		}
	}

	s.started = false
	logging.Info("Health server stopped")
	return nil
}

// SetHealthy sets the healthy status
func (s *Server) SetHealthy(healthy bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.healthy = healthy
}

// SetReady sets the ready status
func (s *Server) SetReady(ready bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ready = ready
}

// healthHandler handles /health endpoint
func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	healthy := s.healthy
	s.mu.RUnlock()

	if healthy {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("UNHEALTHY"))
	}
}

// readyHandler handles /ready endpoint
func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	ready := s.ready
	s.mu.RUnlock()

	if ready {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("READY"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("NOT READY"))
	}
}

// liveHandler handles /live endpoint (liveness probe)
func (s *Server) liveHandler(w http.ResponseWriter, r *http.Request) {
	// Liveness probe - always returns OK if server is running
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ALIVE"))
}
