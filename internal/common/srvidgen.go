package common

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
)

// GenerateServerID creates a unique identifier for a server based on its attributes.
// This function takes the raw URL of the server as input and returns a SHA-256 hash
// of the URL and port as a hexadecimal string. This ensures a unique ID for each server,
// Note: Each server has a unique URL and port combination.
//
// Example usage:
// serverURL := "http://127.0.0.1:5000"
// id := GenerateServerID(serverURL)
// fmt.Println("Server ID:", id)
func GenerateServerID(serverURL string) (string, error) {
	// Parse the raw URL to extract components.
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse server URL: %v", err)
	}

	// Concatenate the URLs hostname and port to form a unique string for hashing.
	hostName := parsedURL.Hostname()
	port := parsedURL.Port()

	if port == "" {
		return "", fmt.Errorf("failed to parse server port: %s", port)
	}

	uniqueString := hostName + ":" + port

	// Compute the SHA-256 hash of the unique string.
	hash := sha256.Sum256([]byte(uniqueString))

	// Encode the hash to a hexadecimal string for easier handling and return.
	hexString := hex.EncodeToString(hash[:])

	return hexString, nil
}

// nolint:lll
/*
1. UUID (Universally Unique Identifier):
Pros: Easy to generate using Go's standard library (uuid package) or third-party libraries. UUIDs are highly unlikely to collide, making them ideal for ensuring uniqueness across a distributed system.
Cons: UUIDs can be long and unwieldy, making them less user-friendly for manual administration tasks. They also don't convey any information about the server they represent.

2. Database-Generated IDs:
Pros: If the load balancer integrates with a database, I can leverage auto-incrementing IDs or other database mechanisms to ensure uniqueness. This approach is straightforward and allows for easy ordering.
Cons: This approach introduces a dependency on external systems (the database), which could be a point of failure or a bottleneck.

3. Combination of Hostname/IP and Port:
Pros: This method is simple and inherently meaningful, as it directly relates to how the server is accessed. It's also relatively easy to implement and understand.
Cons: I might run into uniqueness issues if your system reuses hostnames/IPs with different ports for distinct servers, or if my environment allows for dynamic IP allocation.

4. Hashing Server Attributes (***)
Pros: I can generate a hash (e.g., SHA-256) of several server attributes (RawURL, port, ip and other distinguishing features) to create a unique ID. This method can provide a good balance between meaningfulness and uniqueness.
Cons: There's a small chance of hash collisions, although this is very unlikely with strong hash functions. Also, hashing might require more computation.

5. Centralized ID Service:
Pros: Services like ZooKeeper, etcd, or Consul can be used to generate and manage unique IDs in a distributed system. They offer features like leader election and distributed locks, which can help in managing IDs across a cluster.
Cons: Introduces external dependencies and complexity into your system. Requires management and monitoring of the ID service itself.
*/
