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

	// Find the minimum number of connections.
	minConns := int(^uint(0) >> 1)
	var candidates []Server

	// Step 1: Identify servers with the lowest active connections.
	for _, srv := range servers {
		if srv.IsAlive() {
			conn := srv.GetActiveConnections()
			if conn < minConns {
				minConns = conn
				candidates = []Server{srv} // Start a new list with this server
			} else if conn == minConns {
				candidates = append(candidates, srv) // Add to the list of candidates
			}
		}
	}

	// Step 2: If only one server has the least connections, return it.
	if len(candidates) == 1 {
		return candidates[0]
	}

	// Step 3: If multiple servers have the least connections, use Round-Robin.
	if len(candidates) > 1 {
		// Increment lastSelectedIndex safely.
		lc.lastSelectedIndex = (lc.lastSelectedIndex + 1) % len(servers)

		// Find the next server in the candidates slice that matches the index in the full server list.
		for lc.lastSelectedIndex < len(servers) {
			for _, candidate := range candidates {
				if servers[lc.lastSelectedIndex] == candidate {
					return candidate
				}
			}

			lc.lastSelectedIndex = (lc.lastSelectedIndex + 1) % len(servers)
		}
	}

	// Step 4: If no servers are alive or available, return nil.
	return nil
}
