package utils

import (
	"log"
	"os"
	"strconv"
)

// GetPorts retrieves the starting and load balancer port numbers from environment variables.
// If environment variables are not set or contain invalid values, it defaults to predefined ports.
// STARTING_PORT is the first port for n number of server instances, and LOAD_BALANCER_PORT is for the load balancer.
// Defaults: startingPort :8000 and loadBalancerPort :8080
func GetPorts() (int, int) {
	startingPort := 8000     // Default starting port
	loadBalancerPort := 8080 // Default load balancer port

	if sp, exists := os.LookupEnv("STARTING_PORT"); exists {
		if p, err := strconv.Atoi(sp); err == nil {
			startingPort = p
		} else {
			log.Printf("invalid STARTING_PORT value, using default %d\n", startingPort)
		}
	}

	if lp, exists := os.LookupEnv("LOAD_BALANCER_PORT"); exists {
		if p, err := strconv.Atoi(lp); err == nil {
			loadBalancerPort = p
		} else {
			log.Printf("invalid LOAD_BALANCER_PORT value, using default %d\n", loadBalancerPort)
		}
	}

	return startingPort, loadBalancerPort
}
