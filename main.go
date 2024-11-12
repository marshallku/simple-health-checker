package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marshallku/statusy/config"
	health "github.com/marshallku/statusy/pkg"
	"github.com/marshallku/statusy/server"
	"github.com/marshallku/statusy/store"
)

func main() {
	mode := flag.String("mode", "cli", "Mode to run in (server or health)")
	configFile := flag.String("config", "config.yaml", "Path to configuration file")

	flag.Parse()

	var err error
	var cfg *config.Config
	cfg, err = config.LoadConfig(*configFile)

	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if *mode != "server" {
		health.Check(cfg, nil)
	} else {
		store := store.NewStore()
		server := server.NewServer(store)

		go func() {
			if err := server.Start(); err != nil {
				log.Fatal(err)
			}
		}()

		ticker := time.NewTicker(time.Duration(cfg.CheckInterval) * time.Second)
		defer ticker.Stop()

		for {
			health.Check(cfg, store)
			<-ticker.C
		}
	}
}
