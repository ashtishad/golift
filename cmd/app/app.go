package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ashtishad/golift/internal/domain"
)

// StartServers launches n number of HTTP servers and returns them for management.
func StartServers(startingPort int, n int) []*http.Server {
	var servers []*http.Server

	for i := 0; i < n; i++ {
		server := &http.Server{
			Addr: fmt.Sprintf(":%d", startingPort+i),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello World from server on port %d!", startingPort+i)
			}),
		}

		servers = append(servers, server)

		go func(s *http.Server) {
			if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Failed to start server: %v", err)
			}
		}(server)

		log.Printf("Server listening on port %s", server.Addr)
	}

	return servers
}

func StartLoadBalancer(loadBalancerPort int, startingPort int, srvCnt int) {
	lc := domain.LeastConnection{}
	serverPool := domain.NewServerPool(&lc)

	for i := 0; i < srvCnt; i++ {
		port := startingPort + i
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

	log.Printf("Load Balancer listening on port %d", loadBalancerPort)

	_ = http.ListenAndServe(fmt.Sprintf(":%d", loadBalancerPort), nil)
}
