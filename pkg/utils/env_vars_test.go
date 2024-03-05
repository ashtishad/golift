package utils

import (
	"os"
	"testing"
)

func TestGetPorts(t *testing.T) {
	tests := []struct {
		name                     string
		startingPortEnv          string // Value for STARTING_PORT environment variable
		loadBalancerPortEnv      string // Value for LOAD_BALANCER_PORT environment variable
		expectedStartingPort     int
		expectedLoadBalancerPort int
	}{
		{
			name:                     "Valid environment variables",
			startingPortEnv:          "9000",
			loadBalancerPortEnv:      "9001",
			expectedStartingPort:     9000,
			expectedLoadBalancerPort: 9001,
		},
		{
			name:                     "Invalid STARTING_PORT",
			startingPortEnv:          "invalid",
			loadBalancerPortEnv:      "9001",
			expectedStartingPort:     8000, // Defaults due to invalid STARTING_PORT
			expectedLoadBalancerPort: 9001,
		},
		{
			name:                     "Invalid LOAD_BALANCER_PORT",
			startingPortEnv:          "9000",
			loadBalancerPortEnv:      "invalid",
			expectedStartingPort:     9000,
			expectedLoadBalancerPort: 8080, // Defaults due to invalid LOAD_BALANCER_PORT
		},
		{
			name:                     "Environment variables not set",
			expectedStartingPort:     8000, // Default values
			expectedLoadBalancerPort: 8080,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup cleanup once at the beginning of each test case
			t.Cleanup(func() {
				_ = os.Unsetenv("STARTING_PORT")
				_ = os.Unsetenv("LOAD_BALANCER_PORT")
			})

			// Set environment variables as per the test case
			if tt.startingPortEnv != "" {
				_ = os.Setenv("STARTING_PORT", tt.startingPortEnv)
			}
			if tt.loadBalancerPortEnv != "" {
				_ = os.Setenv("LOAD_BALANCER_PORT", tt.loadBalancerPortEnv)
			}

			startingPort, loadBalancerPort := GetPorts()

			if startingPort != tt.expectedStartingPort || loadBalancerPort != tt.expectedLoadBalancerPort {
				t.Errorf("GetPorts() = (%d, %d), want (%d, %d)", startingPort, loadBalancerPort, tt.expectedStartingPort, tt.expectedLoadBalancerPort)
			}
		})
	}
}
