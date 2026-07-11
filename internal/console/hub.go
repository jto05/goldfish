package console

import (
	// "log"
	"sync"
)

/*
Hub
The hub that tracks a list of clients, a channel for registering clients,
a channel of for unregistering clients, and a channel for broadcasting;
uses a "chan string" as a unique identifier for a client. Clients
is a map for easier deletion
*/
type Hub struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan string

	mu sync.RWMutex
}

func NewHub() (*Hub, error) {
	return nil, nil
}

func (*Hub) Run() error {
	return nil
}

// TODO: define Hub struct with ring buffer and a map of subscriber channels
// TODO: write NewHub(bufSize int) *Hub
// TODO: write Run() — single goroutine selecting on broadcast/register/unregister
// TODO: write Broadcast(line string) — send a line to all subscribers
