package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jto05/goldfish/internal/api"
	"github.com/jto05/goldfish/internal/config"
	"github.com/jto05/goldfish/internal/process"
)

func main() {
	// load config
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg) // print contents for logging

	// build process manager
	m, err := process.NewManager(cfg.Server)
	if err != nil {
		log.Fatal(err)
	}

	// start Minecraft server
	err = m.Start()
	if err != nil {
		log.Fatal(err)
	}

	// start http/websocket server
	h := api.NewHandler(m, m.Hub())
	go func() {
		log.Printf("listening on %s", cfg.API.ListenAddr)
		log.Fatal(http.ListenAndServe(cfg.API.ListenAddr, h.Routes()))
	}()

	// wait for CTRL+C input to stop dashboard
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// stop Minecraft server when exited
	err = m.Stop()
	if err != nil {
		log.Println(err)
	}
}
