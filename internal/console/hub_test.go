package console

import (
	"testing"
	"time"
)

// waitForLine blocks until a string is received on ch or the test times out after 1 second.
func waitForLine(t *testing.T, ch chan string) string {
	t.Helper()
	select {
	case line := <-ch:
		return line
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for message")
		return ""
	}
}

// TestHub_Register verifies that a registered client is tracked in the clients map
// and receives broadcasted messages.
func TestHub_Register(t *testing.T) {
	hub, _ := NewHub(10)
	go hub.Run()

	client := make(chan string, 1)
	hub.register <- client

	// assert client is in the clients map
	if !hub.clients[client] {
		t.Error("expected client to be in clients map after registering")
	}

	// assert client receives a broadcast
	hub.broadcast <- "hello"
	if line := waitForLine(t, client); line != "hello" {
		t.Errorf("expected %q, got %q", "hello", line)
	}
}

// TestHub_Unregister verifies that an unregistered client is removed from the clients map
// and its channel is closed.
func TestHub_Unregister(t *testing.T) {
	hub, _ := NewHub(10)
	go hub.Run()

	client := make(chan string, 1)
	hub.register <- client
	hub.unregister <- client

	// fix race condition
	time.Sleep(10 * time.Millisecond)

	// assert client is removed from clients map
	if hub.clients[client] {
		t.Error("expected client to be removed from clients map after unregistering")
	}

	// assert client channel is closed
	_, open := <-client
	if open {
		t.Error("expected client channel to be closed after unregistering")
	}
}

// TestHub_MultipleClients verifies that a broadcast is received by all registered clients.
func TestHub_MultipleClients(t *testing.T) {
	hub, _ := NewHub(10)
	go hub.Run()

	client1 := make(chan string, 1)
	client2 := make(chan string, 1)
	hub.register <- client1
	hub.register <- client2

	hub.broadcast <- "hello"

	if line := waitForLine(t, client1); line != "hello" {
		t.Errorf("client1: expected %q, got %q", "hello", line)
	}
	if line := waitForLine(t, client2); line != "hello" {
		t.Errorf("client2: expected %q, got %q", "hello", line)
	}
}

// TestHub_RingBufferCatchUp verifies that a newly registered client receives
// historical lines from the ring buffer in order before any live messages.
func TestHub_RingBufferCatchUp(t *testing.T) {
	hub, _ := NewHub(10)
	go hub.Run()

	hub.broadcast <- "line1"
	hub.broadcast <- "line2"
	hub.broadcast <- "line3"

	// wait for broadcasts to be processed before registering
	time.Sleep(10 * time.Millisecond)

	client := make(chan string, 3)
	hub.register <- client

	if line := waitForLine(t, client); line != "line1" {
		t.Errorf("expected %q, got %q", "line1", line)
	}
	if line := waitForLine(t, client); line != "line2" {
		t.Errorf("expected %q, got %q", "line2", line)
	}
	if line := waitForLine(t, client); line != "line3" {
		t.Errorf("expected %q, got %q", "line3", line)
	}
}
