package console

import (
	// "log"
	"log"
)

// Hub that tracks a list of clients, a channel for registering clients,
// a channel of for unregistering clients, and a channel for broadcasting;
// uses a "chan string" as a unique identifier for a client. Clients
// is a map for easier deletion
type Hub struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan string
	history    *RingBuffer[string]
}

// NewHub allocates a Hub with a ring buffer of the given size for console history.
func NewHub(bufferSize int) (*Hub, error) {
	return &Hub{
		clients:    make(map[chan string]bool),
		register:   make(chan chan string),
		unregister: make(chan chan string),
		broadcast:  make(chan string),
		history:    NewRingBuffer[string](bufferSize),
	}, nil // TODO: error handling?
}

// Run starts the hub's event loop and must be called in its own goroutine.
// It selects on register, unregister, and broadcast channels.
func (h *Hub) Run() error {
	for {
		select {
		// register a session and add to list
		case client := <-h.register:

			// register to client list
			h.clients[client] = true

			// send hub history to client
			for _, line := range h.history.ToSlice() {
				client <- line
			}

		// unregister a session and remove from list
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client)

		// broadcast messages to all sessions
		case msg := <-h.broadcast:

			for client := range h.clients {
				client <- msg
			}
			h.history.Push(msg) // push to history
			log.Printf("[Broadcast]: %s", msg)
		}
	}
}

// Broadcast sends a line to the hub's broadcast channel to be fanned out to all clients.
func (h *Hub) Broadcast(line string) {
	h.broadcast <- line
}
