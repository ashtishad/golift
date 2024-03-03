package domain

// ServerPool manages a dynamic set of server instances for load balancing.
type ServerPool interface {
	// AddServer adds a new server to the pool, handling errors like duplicates.
	AddServer(server Server) error

	// RemoveServer removes a server by ID, useful for maintenance or decommissioning.
	RemoveServer(serverID string) error

	// GetServer retrieves a server by ID for status checks or updates.
	GetServer(serverID string) (Server, error)

	// ListServers lists all servers, aiding in monitoring and scaling decisions.
	ListServers() []Server

	// SelectServer picks a server based on the load balancing strategy, optimizing request distribution.
	SelectServer() Server

	// UpdateServerStatus changes a server's alive status, allowing for dynamic health management.
	UpdateServerStatus(serverID string, alive bool) error
}
