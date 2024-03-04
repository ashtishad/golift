package domain

// LoadBalancer interface defines the method for selecting a server from a list.
// It abstracts the strategy used to distribute incoming requests among available servers,
// enabling the implementation of various load balancing algorithms (e.g., Round Robin, Least Connections).
type LoadBalancer interface {
	SelectServer(servers []Server) Server
}
