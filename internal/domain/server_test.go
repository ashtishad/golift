package domain

import (
	"testing"
)

// TestServerAliveStatus tests the SetAlive and IsAlive methods.
func TestServerAliveStatus(t *testing.T) {
	// Backend Server Setup
	serverURL := "http://127.0.0.1:5000"

	srv1, err := NewServer(serverURL)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}

	// Test setting alive status
	srv1.SetAlive(true)
	if !srv1.IsAlive() {
		t.Errorf("expected server to be alive, got not alive")
	}

	srv1.SetAlive(false)
	if srv1.IsAlive() {
		t.Errorf("expected server to be not alive, got alive")
	}
}

// TestNewServerCreation tests the NewServer function for successful server creation.
func TestNewServerCreation(t *testing.T) {
	serverURL := "http://127.0.0.1:5001"

	srv2, err := NewServer(serverURL)
	if err != nil {
		t.Fatalf("NewServer() error = %v, wantErr %v", err, false)
	}

	// Type assertion to access the unexported fields for testing
	s, ok := srv2.(*server)
	if !ok {
		t.Fatalf("Failed to assert the type of server to *server")
	}

	// Check if the server URL is correctly set
	if s.url.String() != serverURL {
		t.Errorf("NewServer() URL = %v, want %v", s.url, serverURL)
	}

	// Check if the server is initialized with alive = true
	if !s.alive {
		t.Errorf("NewServer() server should be initialized as alive")
	}
}
