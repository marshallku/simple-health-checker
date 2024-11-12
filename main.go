package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marshallku/statusy/config"
	health "github.com/marshallku/statusy/pkg"
)

func main() {
	configFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	var err error
	var cfg *config.Config
	cfg, err = config.LoadConfig(*configFile)

	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	health.Check(cfg)
}
