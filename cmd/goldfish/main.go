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
}
