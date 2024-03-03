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
	// SetAlive updates the server's alive status, enabling dynamic health management.
	SetAlive(alive bool)

	// IsAlive reports the current alive status, facilitating health checks and load balancing decisions.
	IsAlive() bool

	// GetURL provides the server's URL, essential for routing and proxying requests.
	GetURL() *url.URL

	// GetActiveConnections returns the current count of active connections, aiding in load balancing strategies like Least Connections.
	GetActiveConnections() int

	// Serve proxies an incoming HTTP request to the server, central to the server's request handling capability.
	Serve(w http.ResponseWriter, r *http.Request)
}

// server implements the Server interface, representing a backend server.
type server struct {
	url          *url.URL               // URL of the server.
	alive        bool                   // Indicates whether the server is alive.
	mux          sync.RWMutex           // Protects access to the server's state.
	activeCons   int32                  // Count of active connections, managed atomically.
	reverseProxy *httputil.ReverseProxy // Used to forward requests to the server.
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

// NewServer creates a new server instance with the specified URL and reverse proxy.
func NewServer(rawURL string) (Server, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse rawURL:%w", err)
	}

	return &server{
		url:          parsedURL,
		alive:        true, // Will use health checks to update.
		activeCons:   0,
		reverseProxy: httputil.NewSingleHostReverseProxy(parsedURL),
	}, nil
}
