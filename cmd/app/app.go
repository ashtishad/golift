package app

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/ashtishad/golift/internal/domain"
)

// StartServers launches n number of HTTP servers and returns them for management.
func StartServers(startingPort int, n int) []*http.Server {
	servers := make([]*http.Server, 0, n)

	for i := range n {
		srv := &http.Server{
			Addr: net.JoinHostPort("", strconv.Itoa(startingPort+i)),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				_, _ = fmt.Fprintf(w, "Hello World from server on port %d!", startingPort+i)
			}),
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       15 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		}

		servers = append(servers, srv)

		go func(s *http.Server) {
			if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("failed to start server: %v", err)
			}
		}(srv)

		log.Printf("Server-%d listening on port %s", i+1, srv.Addr)
	}

	return servers
}

func StartLoadBalancer(loadBalancerPort, port, srvCnt int) {
	lc := domain.LeastConnection{}
	serverPool := domain.NewServerPool(&lc, srvCnt)

	for range srvCnt {
		port++
		serverURL := fmt.Sprintf("http://localhost:%d", port)
		srv, err := domain.NewServer(serverURL)
		if err != nil {
			log.Fatalf("error creating server instance for URL '%s': %v", serverURL, err)
		}

		if err := serverPool.AddServer(srv); err != nil {
			return
		}
	}

	// Setup and start the load balancer HTTP server.
	http.HandleFunc("/", proxyRequestHandler(serverPool))

	// Create a custom http.Server with timeouts
	s := &http.Server{
		Addr:         net.JoinHostPort("", strconv.Itoa(loadBalancerPort)),
		Handler:      proxyRequestHandler(serverPool),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	log.Printf("Load Balancer listening on port %d", loadBalancerPort)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start load balancer: %v", err)
	}
}
