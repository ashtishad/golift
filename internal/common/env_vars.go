package common

import (
	"log/slog"
	"os"
	"strconv"
)

// Config holds environment-based configuration parameters
type Config struct {
	APIHost          string
	NumOfServers     int
	StartingPort     int
	LoadBalancerPort string
}

// LoadConfig loads the configuration with defaults.
func LoadConfig(l *slog.Logger) *Config {
	config := Config{
		APIHost:          "127.0.0.1",
		NumOfServers:     5,
		StartingPort:     8000,
		LoadBalancerPort: "8080",
	}

	requiredVars := []string{"API_HOST", "STARTING_PORT", "LOAD_BALANCER_PORT", "NUM_OF_SERVERS"}

	for _, varName := range requiredVars {
		value := os.Getenv(varName)
		if value != "" {
			switch varName {
			case "API_HOST":
				config.APIHost = value
			case "STARTING_PORT":
				port, err := strconv.Atoi(value)
				if err != nil {
					l.Error("invalid STARTING_PORT", "err", err, "port", port) // Log error & continue
					continue
				}
				config.StartingPort = port
			case "LOAD_BALANCER_PORT":
				config.LoadBalancerPort = value
			case "NUM_OF_SERVERS":
				num, err := strconv.Atoi(value)
				if err != nil {
					l.Error("invalid NUM_OF_SERVERs", "err", err, "srv_cnt", num)
					continue
				}
				config.NumOfServers = num
			}
		} else {
			l.Warn("environment variable is not defined. Using default", "varName", varName)
		}
	}

	return &config
}
