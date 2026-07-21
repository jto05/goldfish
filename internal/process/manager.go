package process

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/jto05/goldfish/internal/config"
	"github.com/jto05/goldfish/internal/console"
)

type ServerStatus string

const (
	Stopped     ServerStatus = "stopped"
	Running     ServerStatus = "running"
	BUFFER_SIZE int          = 500
)

type Manager struct {
	cfg       config.ServerConfig
	status    ServerStatus
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	startedAt time.Time
	hub       *console.Hub
	mu        sync.Mutex
}

func NewManager(cfg config.ServerConfig) (*Manager, error) {
	// build manager
	mng := Manager{
		cfg:    cfg,
		status: Stopped,
		// cmd and stdin are initalized when buildCmd() called
		startedAt: time.Time{},
		hub:       console.NewHub(BUFFER_SIZE),
	}
	return &mng, nil
}

func (m *Manager) buildCmd() (io.ReadCloser, error) {
	// parse cmd fields
	parts := strings.Fields(m.cfg.StartCmd)
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = m.cfg.WorkDir // set workign directory for command

	// set stdin pipe
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	m.stdin = stdin

	// get stdout pipe
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	m.cmd = cmd

	return stdout, nil
}

func (m *Manager) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status == Running {
		return fmt.Errorf("server already running")
	}

	// get stdout and build cmd
	stdout, err := m.buildCmd()
	if err != nil {
		return err
	}

	err = m.cmd.Start()
	if err != nil {
		return err
	}

	m.status = Running
	m.startedAt = time.Now()

	// start hub
	go m.hub.Run()

	// pipe stdout to hub
	go func() {
		// using scanner  for cleaner line by line outputs
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m.hub.Broadcast(scanner.Text()) // send each line through hub
		}
	}()

	// update Manager status
	go func() {
		// wait for cmd to stop
		m.cmd.Wait()
		m.mu.Lock()
		m.status = Stopped
		m.mu.Unlock()
	}()

	return nil
}

func (m *Manager) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.status == Stopped {
		return fmt.Errorf("server is already stopped")
	}

	// send stop command to stdin of command
	_, err := fmt.Fprintln(m.stdin, "stop") // NOTE: hardcoded "stop"?
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Status() ServerStatus {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.status
}

func (m *Manager) SendCommand(cmd string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, err := fmt.Fprintln(m.stdin, cmd)
	if err != nil {
		return err
	}
	return err
}
