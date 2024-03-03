package domain

import "sync"

// LoadBalancer interface defines the method for selecting a server from a list.
// It abstracts the strategy used to distribute incoming requests among available servers,
// enabling the implementation of various load balancing algorithms (e.g., Round Robin, Least Connections).
type LoadBalancer interface {
	SelectServer(servers []Server) Server
}

type LeastConnection struct {
	lastSelectedIndex int
	mux               sync.Mutex
}

// SelectServer selects a server based on the least connections strategy with a Round-Robin tiebreaker.
// 1: Select servers with the lowest active connections.
// 2: Directly assign the request to a lone server with the fewest connections.
// 3: If multiple servers share the lowest count, employ Round Robin to assign the request.
// 4. If no servers are alive or available, return nil.
func (lc *LeastConnection) SelectServer(servers []Server) Server {
	lc.mux.Lock()
	defer lc.mux.Unlock()

	// Step 1: Identify Servers with the lowest active connections.
	minConns := int(^uint(0) >> 1) // Initialize with the maximum int value.
	candidates := make([]Server, 0)

	for _, srv := range servers {
		if srv.IsAlive() {
			conn := srv.GetActiveConnections()
			if conn < minConns {
				minConns = conn
				candidates = candidates[:0] // Reset the slice
				candidates = append(candidates, srv)
			} else if conn == minConns {
				candidates = append(candidates, srv)
			}
		}
	}

	// Step 2: Single Server Assignment
	if len(candidates) == 1 {
		return candidates[0]
	}

	// Step 3: Round-Robin Tiebreaker
	if len(candidates) > 1 {
		lc.lastSelectedIndex = (lc.lastSelectedIndex + 1) % len(candidates)
		return candidates[lc.lastSelectedIndex]
	}

	// If no servers are alive or available, return nil.
	return nil
}
