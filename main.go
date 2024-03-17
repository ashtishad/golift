package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashtishad/golift/internal/common"
	"github.com/ashtishad/golift/internal/domain"
	"github.com/ashtishad/golift/internal/transport"
)

func main() {
	// init slogger
	handlerOpts := common.GetSlogConf()
	logger := slog.New(slog.NewTextHandler(os.Stdout, handlerOpts))
	slog.SetDefault(logger)

	// load config
	conf := common.LoadConfig(logger)

	// Start servers and load balancer.
	servers := startServers(conf, logger)
	go startLoadBalancer(conf, logger)

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
func startServers(conf *common.Config, l *slog.Logger) []*http.Server {
	startingPort := conf.StartingPort
	n := conf.NumOfServers

	servers := make([]*http.Server, 0, n)

	for i := 0; i < n; i++ {
		server := &http.Server{
			Addr: fmt.Sprintf(":%d", startingPort+i),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = fmt.Fprintf(w, "Hello World from server on port %d!", startingPort+i)
			}),
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       15 * time.Second,
			ReadHeaderTimeout: 2 * time.Second,
		}

		servers = append(servers, server)

		go func(s *http.Server) {
			if err := s.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				l.Error("failed to start the server", "err", err)
				os.Exit(1)
			}
		}(server)

		l.Info(fmt.Sprintf("Server-%d listening at", i+1), "addr", server.Addr)
	}

	return servers
}

func startLoadBalancer(conf *common.Config, l *slog.Logger) {
	loadBalancerPort := conf.LoadBalancerPort
	startingPort := conf.StartingPort
	srvCnt := conf.NumOfServers

	lc := domain.LeastConnection{}
	serverPool := domain.NewServerPool(&lc, srvCnt, l)

	for i := 0; i < srvCnt; i++ {
		port := startingPort + i
		serverURL := fmt.Sprintf("http://localhost:%d", port)
		srv, err := domain.NewServer(serverURL)

		if err != nil {
			l.Error("error creating server instances", "url", serverURL, "err", err)
		}

		if err := serverPool.AddServer(srv); err != nil {
			return
		}
	}

	// Setup and start the load balancer HTTP server.
	handler := transport.ProxyRequestHandler(serverPool, l)
	http.HandleFunc("/", handler)

	// Create a custom http.Server with timeouts
	s := &http.Server{
		Addr:         net.JoinHostPort(conf.APIHost, loadBalancerPort),
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	l.Info("Load balancer listening at", "addr", s.Addr)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		l.Error("failed to start load balancer", "err", err)
	}
}
