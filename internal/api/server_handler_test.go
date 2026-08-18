package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jto05/goldfish/internal/config"
	"github.com/jto05/goldfish/internal/process"
)

func newTestHandler(t *testing.T) *Handler {
	t.Helper()
	m, err := process.NewManager(config.ServerConfig{
		StartCmd: "echo fake",
		WorkDir:  t.TempDir(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return NewHandler(m, m.Hub())
}

func TestServerStatus_Stopped(t *testing.T) {
	h := newTestHandler(t)

	r := httptest.NewRequest("GET", "/api/server/status", nil)
	w := httptest.NewRecorder()
	h.serverStatus(w, r)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", res.StatusCode)
	}

	var body map[string]string
	json.NewDecoder(res.Body).Decode(&body)
	if body["status"] != "stopped" {
		t.Errorf("expected status %q, got %q", "stopped", body["status"])
	}
}

func TestServerStop_WhenStopped(t *testing.T) {
	h := newTestHandler(t)

	r := httptest.NewRequest("POST", "/api/server/stop", nil)
	w := httptest.NewRecorder()
	h.serverStop(w, r)

	if w.Result().StatusCode != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Result().StatusCode)
	}
}

func TestServerCommand_InvalidBody(t *testing.T) {
	h := newTestHandler(t)

	r := httptest.NewRequest("POST", "/api/server/command", strings.NewReader("not json"))
	w := httptest.NewRecorder()
	h.serverCommand(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}

func TestServerCommand_EmptyCommand(t *testing.T) {
	h := newTestHandler(t)

	r := httptest.NewRequest("POST", "/api/server/command", strings.NewReader(`{"command":""}`))
	w := httptest.NewRecorder()
	h.serverCommand(w, r)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Result().StatusCode)
	}
}
