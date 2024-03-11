package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ashtishad/golift/internal/common"
	"github.com/ashtishad/golift/internal/domain"
	"github.com/ashtishad/golift/internal/transport"
)

func main() {
	startingPort, loadBalancerPort := common.GetPorts()
	srvCnt := common.GetServerCount()

	// Start servers and load balancer.
	servers := startServers(startingPort, srvCnt)
	go startLoadBalancer(loadBalancerPort, startingPort, srvCnt)

	// Setup channel to listen for OS interrupt signals for graceful shutdown.
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGINT, syscall.SIGTERM)

	<-quitChan // Block until a signal is received.

	// Shutdown logic for servers.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, server := range servers {
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("error shutting down server: %v", err)
		}
	}

	log.Println("Servers shut down gracefully.")
}

// startServers launches n number of HTTP servers and returns them for management.
func startServers(startingPort int, n int) []*http.Server {
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

func startLoadBalancer(loadBalancerPort, port, srvCnt int) {
	lc := domain.LeastConnection{}
	serverPool := domain.NewServerPool(&lc, srvCnt)

	for range srvCnt {
		port++
		serverURL := fmt.Sprintf("http://localhost:%d", port)
		srv, err := domain.NewServer(serverURL)
		if err != nil {
			log.Fatalf("error creating server instance for URL '%s': %v", serverURL, err)
		}

		if spErr := serverPool.AddServer(srv); spErr != nil {
			return
		}
	}

	// Setup and start the load balancer HTTP server.
	srv := &http.Server{
		Addr:              net.JoinHostPort(os.Getenv("API_HOST"), strconv.Itoa(loadBalancerPort)),
		Handler:           nil,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       15 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	router := http.NewServeMux()
	http.HandleFunc("GET /", transport.ProxyRequestHandler(serverPool))
	srv.Handler = router

	log.Printf("Load Balancer listening on port %d", loadBalancerPort)

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start load balancer: %v", err)
	}
}
