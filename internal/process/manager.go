package process

// TODO: define Manager struct with mu, status, cmd, stdin, startedAt fields
// TODO: write NewManager(cfg config.ServerConfig) *Manager
// TODO: write Start() — launch server jar with os/exec, capture stdout/stderr, transition Starting -> Running
// TODO: write Stop() — write stop command to stdin, transition Stopping -> Stopped
// TODO: write SendCommand(cmd string) — write arbitrary command to stdin
// TODO: write Status() ServerStatus — thread-safe read of current status
// TODO: goroutine to detect process exit and transition back to Stopped
// TODO: pipe process stdout/stderr into hub.Broadcast inside Start()
