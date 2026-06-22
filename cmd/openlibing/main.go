package main

import (
	"os"

	embedded "github.com/openlibing/openlibing-cli/embedded"
	"github.com/openlibing/openlibing-cli/internal/api"
	"github.com/openlibing/openlibing-cli/internal/cli"
	"github.com/openlibing/openlibing-cli/internal/config"
	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/internal/registry"
)

func main() {
	// Load configuration
	cfg, _ := config.LoadConfig()
	auth, _ := config.LoadAuth()

	// Initialize API client
	client := api.NewClient(cfg, auth)

	// Initialize Registry and load SPCs (3-layer)
	reg := registry.NewRegistry()
	if err := reg.LoadAll(embedded.SPCs, "spc"); err != nil {
		// Non-fatal: registry errors are logged but CLI still works
		// (at minimum, 'openlibing list' will show empty)
	}

	// Initialize Engine
	eng := engine.NewEngine(reg, client)

	// Build and execute CLI
	rootCmd := cli.NewRootCmd(eng, reg)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
