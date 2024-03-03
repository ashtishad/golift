package domain

import "sync"

// ServerPooler defines operations for managing a dynamic set of server instances for load balancing.
type ServerPooler interface {
	// AddServer adds a new server to the pool, generates id and handling errors like duplicates.
	AddServer(srv Server) error

	// RemoveServer removes a server by ID, useful for maintenance or decommissioning.
	RemoveServer(srvID string) error

	// GetServer retrieves a server by ID for status checks or updates.
	GetServer(srvID string) (Server, error)

	// ListServers lists all servers, aiding in monitoring and scaling decisions.
	ListServers() []Server

	// SelectServer picks a server based on the load balancing strategy, optimizing request distribution.
	SelectServer() Server

	// UpdateServerStatus changes a server's alive status, allowing for dynamic health management.
	UpdateServerStatus(srvID string, alive bool) error
}

type serverPool struct {
	servers map[string]server
	mux     sync.RWMutex
}
