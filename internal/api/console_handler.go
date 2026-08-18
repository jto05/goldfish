package api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to Websocket
// connections
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// TODO: implement meaningful error messages
func (h *Handler) serverConsole(w http.ResponseWriter, r *http.Request) {
	// upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// register/unregister client
	client := make(chan string, 64)
	h.hub.Register(client)
	defer h.hub.Unregister(client)

	// block until client disconnects
	for line := range client {
		err := conn.WriteMessage(websocket.TextMessage, []byte(line))
		if err != nil {
			break
		}
	}
}
