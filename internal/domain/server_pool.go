package domain

// ServerPool defines the operations for managing a dynamic pool of server instances.
type ServerPool interface {
	AddServer(srv Server) error                     // Add a new server to the pool.
	RemoveServer(id string) error                   // Remove a server from the pool by its identifier.
	GetServer(id string) (Server, error)            // Retrieve a server by its identifier.
	ListServers() []Server                          // List all servers in the pool.
	SelectServer() Server                           // Select a server based on the load balancing strategy (e.g., least connections, round robin).
	UpdateServerStatus(id string, alive bool) error // Update the alive status of a server.
}
