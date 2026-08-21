package api

import (
	"net/http"

	"github.com/jto05/goldfish/internal/console"
	"github.com/jto05/goldfish/internal/process"
)

type Handler struct {
	manager *process.Manager
	hub     *console.Hub
}

func NewHandler(manager *process.Manager, hub *console.Hub) *Handler {
	return &Handler{manager, hub}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/server/status", h.serverStatus)
	mux.HandleFunc("POST /api/server/start", h.serverStart)
	mux.HandleFunc("POST /api/server/stop", h.serverStop)
	mux.HandleFunc("POST /api/server/command", h.serverCommand)

	mux.HandleFunc("GET /ws/console", h.serverConsole)
	return cors(mux)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
