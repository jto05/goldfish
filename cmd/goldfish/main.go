package main

import (
	"fmt"
	"log"

	"github.com/jto05/goldfish/internal/config"
)

func main() {
	cfg, err := config.Load("config.example.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg)

	// TODO: start the hub
	// TODO: start the manager
	// TODO: print console output to terminal
	// TODO: handle Ctrl-C (OS signal) to call Stop() cleanly
}
