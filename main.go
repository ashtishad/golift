package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashtishad/golift/cmd/app"
	"github.com/ashtishad/golift/pkg/utils"
)

func main() {
	startingPort, loadBalancerPort := utils.GetPorts()
	srvCnt := utils.GetServerCount()

	// Start servers and load balancer.
	servers := app.StartServers(startingPort, srvCnt)
	go app.StartLoadBalancer(loadBalancerPort, startingPort, srvCnt)

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
