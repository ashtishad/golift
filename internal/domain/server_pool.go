package domain

import (
	"fmt"
	"sync"

	"github.com/ashtishad/golift/pkg/srvidgen"
)

// ServerPooler defines operations for managing a dynamic set of server instances for load balancing.
type ServerPooler interface {
	// AddServer adds a new server to the pool, generates server id and handling errors like duplicates.
	AddServer(srv Server) error

	// RemoveServer removes a server by ID, useful for maintenance or decommissioning.
	RemoveServer(srvID string) error

	// GetServer retrieves a server by ID for status checks or updates.
	GetServer(srvID string) (Server, error)

	// ListServers lists all servers, aiding in monitoring and scaling decisions.
	ListServers() []Server

	// SelectServer picks a server based on the load balancing strategy.
	SelectServer() Server

	// UpdateServerStatus changes a server's alive status, allowing for dynamic health management.
	UpdateServerStatus(srvID string, alive bool) error
}

type serverPool struct {
	servers  map[string]Server
	mux      sync.RWMutex
	strategy LoadBalancer // For selecting a server based on Load Balancing Strategy
}

// AddServer adds a new server to the pool, generates server id and handling errors like duplicates.
// returns an error if something went wrong.
func (sp *serverPool) AddServer(srv Server) error {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	// Generate server id with hash value of server URL and port
	srvID, err := srvidgen.GenerateServerID(srv.GetURL().String())
	if err != nil {
		return fmt.Errorf("failed to generate server id from server URL: %w", err)
	}

	// Check if the server already exists in the pool to avoid duplicates.
	if _, exists := sp.servers[srvID]; exists {
		return fmt.Errorf("server with ID %s already exists in the pool", srvID)
	}

	srv.SetID(srvID)

	// Add the server to the pool.
	sp.servers[srvID] = srv

	return nil
}

// RemoveServer removes a server by ID, It returns an error if the server to be removed does not exist in the pool.
// Used The `delete` function, safe to call even if the key is not present in the map,
// However, the existence check is performed to provide specific error feedback.
func (sp *serverPool) RemoveServer(srvID string) error {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	// Check if the server exists in the pool.
	if _, exists := sp.servers[srvID]; !exists {
		return fmt.Errorf("server with ID %s does not exist in the pool", srvID)
	}

	// Remove the server from the pool.
	delete(sp.servers, srvID)

	return nil
}

// GetServer retrieves a server by ID for status checks or updates.
// It returns the server if found. Returns an error If the server is not found,
func (sp *serverPool) GetServer(srvID string) (Server, error) {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	// Attempt to retrieve the server from the pool using its ID.
	if srv, exists := sp.servers[srvID]; exists {
		return srv, nil
	}

	// If the server is not found, return an error.
	return nil, fmt.Errorf("server with ID %s not found", srvID)
}

// ListServers lists all servers, returns a slice containing all the servers currently in the pool.
func (sp *serverPool) ListServers() []Server {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	servers := make([]Server, 0)

	// Iterate over the map of servers and add each server to the slice.
	for _, srv := range sp.servers {
		servers = append(servers, srv)
	}

	return servers
}

// UpdateServerStatus changes a server's alive status, allowing for dynamic health management.
func (sp *serverPool) UpdateServerStatus(srvID string, alive bool) error {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	if srv, exists := sp.servers[srvID]; exists {
		// If the server is found, update its alive status.
		srv.SetAlive(alive)
		return nil
	}

	// If the server is not found, return an error.
	return fmt.Errorf("server with ID %s not found", srvID)
}

// SelectServer picks a server based on the load balancing strategy.
func (sp *serverPool) SelectServer() Server {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	var serversSlice []Server
	for _, srv := range sp.servers {
		serversSlice = append(serversSlice, srv)
	}

	return sp.strategy.SelectServer(serversSlice)
}

func NewServerPool(strategy LoadBalancer) ServerPooler {
	return &serverPool{
		servers:  make(map[string]Server),
		strategy: strategy,
	}
}
