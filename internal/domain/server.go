package domain

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

// Server defines the operations necessary for a server within a load-balanced environment.
type Server interface {
	SetAlive(alive bool)                          // Updates the server's alive status.
	IsAlive() bool                                // Reports the current alive status.
	GetURL() *url.URL                             // Provides the server's URL.
	GetActiveConnections() int                    // Returns the current count of active connections.
	Serve(w http.ResponseWriter, r *http.Request) // Proxies an incoming HTTP request.
	GetID() string                                // Returns a unique identifier for the server.
	SetID(srvID string)                           // Sets a unique identifier for the server.
}

// server implements the Server interface, representing a backend server.
type server struct {
	id           string                 // Unique identifier for the server.
	url          *url.URL               // URL of the server.
	alive        bool                   // Indicates whether the server is alive.
	mux          sync.RWMutex           // Protects access to the server's state.
	activeCons   int32                  // Count of active connections, managed atomically.
	reverseProxy *httputil.ReverseProxy // Used to forward requests to the server.
}

// NewServer creates a new server instance with the specified URL and reverse proxy.
func NewServer(rawURL string) (Server, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse rawURL: %w", err)
	}

	return &server{
		url:          parsedURL,
		alive:        true, // Will use health checks to update.
		activeCons:   0,
		reverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}, nil
}

// GetID returns the server's unique identifier.
func (s *server) GetID() string {
	return s.id
}

// SetID returns the server's unique identifier.
func (s *server) SetID(srvID string) {
	s.id = srvID
}

// SetAlive updates the server's alive status. It safely handles concurrent updates.
func (s *server) SetAlive(a bool) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.alive = a
}

// IsAlive returns the current alive status of the server, ensuring thread-safe access.
func (s *server) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()

	return s.alive
}

// GetURL retrieves the server's URL.
func (s *server) GetURL() *url.URL {
	return s.url
}

// GetActiveConnections returns the current number of active connections to the server.
// Using atomic operations for thread-safe access.
func (s *server) GetActiveConnections() int {
	return int(atomic.LoadInt32(&s.activeCons))
}

// Serve forwards the incoming HTTP request to the server using the reverse proxy.
// It increments and decrements the active connection count before and after serving the request.
func (s *server) Serve(rw http.ResponseWriter, req *http.Request) {
	atomic.AddInt32(&s.activeCons, 1)
	defer atomic.AddInt32(&s.activeCons, -1)

	s.reverseProxy.ServeHTTP(rw, req)
}
