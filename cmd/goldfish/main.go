package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jto05/goldfish/internal/config"
	"github.com/jto05/goldfish/internal/process"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg)

	m, err := process.NewManager(cfg.Server)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Start()
	if err != nil {
		log.Fatal(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	if err := m.Stop(); err != nil {
		log.Println(err)
	}
}
