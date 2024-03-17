package domain

import (
	"log/slog"
	"sync"

	"github.com/ashtishad/golift/internal/common"
)

// ServerPooler defines operations for managing a dynamic set of server instances for load balancing.
type ServerPooler interface {
	// AddServer adds a new server to the pool, generates server id and handling errors like duplicates.
	AddServer(srv Server) common.AppError

	// RemoveServer removes a server by ID, useful for maintenance or decommissioning.
	RemoveServer(srvID string) common.AppError

	// GetServer retrieves a server by ID for status checks or updates.
	GetServer(srvID string) (Server, common.AppError)

	// ListServers lists all servers, aiding in monitoring and scaling decisions.
	ListServers() []Server

	// SelectServer picks a server based on the underlying load balancing strategy from LoadBalancer interface.
	SelectServer() Server

	// UpdateServerStatus changes a server's alive status, allowing for dynamic health management.
	UpdateServerStatus(srvID string, alive bool) common.AppError
}

type serverPool struct {
	servers  map[string]Server
	mux      sync.RWMutex
	strategy LoadBalancer
	logger   *slog.Logger
}

// AddServer adds a new server to the pool, generates server id and handling errors like duplicates.
// returns an common.AppError if something went wrong.
func (sp *serverPool) AddServer(srv Server) common.AppError {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	// Generate server id with hash value of server URL and port
	srvID, err := common.GenerateServerID(srv.GetURL().String())
	if err != nil {
		sp.logger.Error("failed to generate server id from server URL", "err", err)
		return common.NewInternalServerError("failed to generate server id from server URL", err)
	}

	// Check if the server already exists in the pool to avoid duplicates.
	if _, exists := sp.servers[srvID]; exists {
		sp.logger.Warn("server ID already exists in the pool", "srv_id", srvID)
		return common.NewConflictError("server id already exists in the pool")
	}

	srv.SetID(srvID)

	// Add the server to the pool.
	sp.servers[srvID] = srv

	return nil
}

// RemoveServer removes a server by ID, It returns an common.AppError if the server to be removed does not exist in the pool.
// Used The `delete` function, safe to call even if the key is not present in the map,
// However, the existence check is performed to provide specific common.AppError feedback.
func (sp *serverPool) RemoveServer(srvID string) common.AppError {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	if _, exists := sp.servers[srvID]; !exists {
		sp.logger.Warn("server ID does not exist in the pool", "srv_id", srvID)
		return common.NewConflictError("server id does not exist in the pool")
	}

	delete(sp.servers, srvID)

	return nil
}

// GetServer retrieves a server by ID for status checks or updates.
// It returns the server if found. Returns an common.AppError If the server is not found,
func (sp *serverPool) GetServer(srvID string) (Server, common.AppError) {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	if srv, exists := sp.servers[srvID]; exists {
		return srv, nil
	}

	// If the server is not found, return a common.AppError.
	sp.logger.Error("server with id not found", "srv_id", srvID)
	return nil, common.NewNotFoundError("server with id not found")
}

// ListServers lists all servers, returns a slice containing all the servers currently in the pool.
func (sp *serverPool) ListServers() []Server {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	servers := make([]Server, 0)

	for _, srv := range sp.servers {
		servers = append(servers, srv)
	}

	return servers
}

// UpdateServerStatus changes a server's alive status, allowing for dynamic health management.
func (sp *serverPool) UpdateServerStatus(srvID string, alive bool) common.AppError {
	sp.mux.Lock()
	defer sp.mux.Unlock()

	if srv, exists := sp.servers[srvID]; exists {
		// If the server is found, update its alive status.
		srv.SetAlive(alive)
		return nil
	}

	// If the server is not found, return a common.AppError.
	sp.logger.Error("server with id not found", "srv_id", srvID)
	return common.NewNotFoundError("server with id not found")
}

// SelectServer picks a server based on the underlying load balancing strategy from LoadBalancer interface.
func (sp *serverPool) SelectServer() Server {
	sp.mux.RLock()
	defer sp.mux.RUnlock()

	serversSlice := make([]Server, 0, len(sp.servers))
	for _, srv := range sp.servers {
		serversSlice = append(serversSlice, srv)
	}

	return sp.strategy.SelectServer(serversSlice)
}

func NewServerPool(strategy LoadBalancer, cnt int, logger *slog.Logger) ServerPooler {
	return &serverPool{
		servers:  make(map[string]Server, cnt),
		strategy: strategy,
		logger:   logger,
	}
}
