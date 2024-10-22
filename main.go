package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marshallku/simple_health_checker/config"
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

	fmt.Printf("Loaded configuration: %+v\n", cfg)

}
