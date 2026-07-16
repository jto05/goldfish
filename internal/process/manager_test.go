package process

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jto05/goldfish/internal/config"
)

// TestMain allows this test binary to act as the fake Minecraft server
func TestMain(m *testing.M) {
	if os.Getenv("GO_FAKE_SERVER") == "1" {
		runFakeServer()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

// runFakeServer simulates a Minecraft server: prints a start message,
// reads commands from stdin, and exits on "stop".
func runFakeServer() {
	fmt.Println("test server started")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("received: %s\n", line)
		if line == "stop" {
			fmt.Println("fake server stopped")
			return
		}
	}
}

// fakeServerCmd returns the start command that launches this test binary
// as a fake server.
func fakeServerCmd() string {
	return fmt.Sprintf("%s -test.run=^$ -test.v", os.Args[0])
}

func fakeCfg() config.ServerConfig {
	return config.ServerConfig{
		WorkDir:  ".",
		StartCmd: fakeServerCmd(),
	}
}

func TestManager_Start(t *testing.T) {
	os.Setenv("GO_FAKE_SERVER", "1")
	defer os.Unsetenv("GO_FAKE_SERVER")

	mgr, err := NewManager(fakeCfg())
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}

	if err := mgr.Start(); err != nil {
		t.Fatalf("Start: %v", err)
	}

	if mgr.Status() != Running {
		t.Errorf("expected status Running, got %s", mgr.Status())
	}
}

func TestManager_DoubleStart(t *testing.T) {
	os.Setenv("GO_FAKE_SERVER", "1")
	defer os.Unsetenv("GO_FAKE_SERVER")

	mgr, _ := NewManager(fakeCfg())
	mgr.Start()

	if err := mgr.Start(); err == nil {
		t.Error("expected error on double start, got nil")
	}
}

func TestManager_Stop(t *testing.T) {
	os.Setenv("GO_FAKE_SERVER", "1")
	defer os.Unsetenv("GO_FAKE_SERVER")

	mgr, _ := NewManager(fakeCfg())
	mgr.Start()

	if err := mgr.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	// give the process time to exit
	time.Sleep(100 * time.Millisecond)

	if mgr.Status() != Stopped {
		t.Errorf("expected status Stopped, got %s", mgr.Status())
	}
}

func TestManager_StopWhenStopped(t *testing.T) {
	mgr, _ := NewManager(fakeCfg())

	if err := mgr.Stop(); err == nil {
		t.Error("expected error stopping an already stopped server, got nil")
	}
}
