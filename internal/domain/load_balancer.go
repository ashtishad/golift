package domain

// LoadBalancer defines the interface for load balancing strategies.
type LoadBalancer interface {
	SelectServer(servers []Server) Server
}
