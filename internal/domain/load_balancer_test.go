package domain

import (
	"net/http"
	"net/url"
	"sync"
	"testing"
)

// MockServer implements the Server interface for testing purposes.
type MockServer struct {
	id                string
	alive             bool
	activeConnections int
	mux               sync.Mutex
}

// TestLeastConnection_SelectServer tests the LeastConnection strategy.
func TestLeastConnection_SelectServer(t *testing.T) {
	mockServers := []*MockServer{
		{id: "server1", alive: true, activeConnections: 1},
		{id: "server2", alive: true, activeConnections: 1},
		{id: "server3", alive: true, activeConnections: 1},
		{id: "server4", alive: true, activeConnections: 1},
		{id: "server5", alive: true, activeConnections: 1},
		{id: "server6", alive: true, activeConnections: 1},
	}

	// Convert mock servers to a slice of Server interfaces
	servers := make([]Server, len(mockServers))
	for i, mockServer := range mockServers {
		servers[i] = mockServer
	}

	// Create an instance of LeastConnection
	lc := &LeastConnection{}

	// Define test scenarios
	testScenarios := []struct {
		name              string
		setup             func()
		expectedServerIDs []string
	}{
		{
			name: "Single Server Selection",
			setup: func() {
				// Only one server with the least connections
				mockServers[0].activeConnections = 4
				mockServers[1].activeConnections = 1
				mockServers[2].activeConnections = 2
				mockServers[3].activeConnections = 0
				mockServers[4].activeConnections = 4
				mockServers[5].activeConnections = 5
			},
			expectedServerIDs: []string{"server4"},
		},
		// {
		// 	name: "Multiple Servers with Equal Connections",
		// 	setup: func() {
		// 		// Multiple servers with equal connections
		// 		for _, server := range mockServers {
		// 			server.activeConnections = 1
		// 		}
		// 	},
		// 	expectedServerIDs: []string{"server1", "server2", "server3", "server4", "server5", "server6"},
		// },
		{
			name: "No Servers Available",
			setup: func() {
				// No servers are alive
				for _, server := range mockServers {
					server.alive = false
				}
			},
			expectedServerIDs: []string{""},
		},
		{
			name: "Multiple Servers with Different Connections",
			setup: func() {
				// Servers with different connections
				mockServers[0].activeConnections = 2
				mockServers[1].activeConnections = 5
				mockServers[2].activeConnections = 3
				mockServers[3].activeConnections = 2
				mockServers[4].activeConnections = 0
				mockServers[5].activeConnections = 0
			},
			expectedServerIDs: []string{"server5", "server6"},
		},
		{
			name: "Round-Robin Cycling",
			setup: func() {
				// Reset the round-robin index
				lc.lastSelectedIndex = -1
				// All servers have the same number of connections
				for _, server := range mockServers {
					server.activeConnections = 1
				}
			},
			expectedServerIDs: []string{"server1", "server2", "server3", "server4", "server5", "server6", "server1"},
		},
	}

	// Run test scenarios
	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Setup the scenario
			scenario.setup()

			// Run the selection process as many times as expectedServerIDs
			for i, expectedID := range scenario.expectedServerIDs {
				selectedServer := lc.SelectServer(servers)
				selectedID := ""
				if selectedServer != nil {
					selectedID = selectedServer.GetID()
				}
				if expectedID != selectedID {
					t.Errorf("Round %d: Expected server ID %s, got %s", i+1, expectedID, selectedID)
				}
			}

			// Reset server states after each scenario
			for _, server := range mockServers {
				server.alive = true
				server.activeConnections = 1
			}
		})
	}
}

func (m *MockServer) Serve(w http.ResponseWriter, r *http.Request) {
	// TODO implement me
	panic("implement me")
}

func (m *MockServer) SetAlive(alive bool) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.alive = alive
}

func (m *MockServer) IsAlive() bool {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.alive
}

func (m *MockServer) GetURL() *url.URL {
	// Simplified for testing; assume all servers have a URL.
	return &url.URL{Scheme: "http", Host: "localhost"}
}

func (m *MockServer) GetActiveConnections() int {
	m.mux.Lock()
	defer m.mux.Unlock()
	return m.activeConnections
}

func (m *MockServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Simplified for testing.
}

func (m *MockServer) GetID() string {
	return m.id
}

func (m *MockServer) SetID(id string) {
	m.id = id
}
