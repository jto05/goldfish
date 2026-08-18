package api

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) serverStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": string(h.manager.Status()),
	})
}

func (h *Handler) serverStop(w http.ResponseWriter, r *http.Request) {
	err := h.manager.Stop()
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // return 409 if error found
		return
	}
	w.WriteHeader(http.StatusNoContent) // return 204, no content
}

func (h *Handler) serverStart(w http.ResponseWriter, r *http.Request) {
	err := h.manager.Start()
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict) // return 409 if error found
		return
	}
	w.WriteHeader(http.StatusNoContent) // return 204, no content
}

func (h *Handler) serverCommand(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Command string `json:"command"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	// check request body
	if err != nil { // invalid body
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if body.Command == "" { // empty body
		http.Error(w, "command is required", http.StatusBadRequest)
		return
	}

	// return 500 internal error if SendComand fails
	err = h.manager.SendCommand(body.Command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
