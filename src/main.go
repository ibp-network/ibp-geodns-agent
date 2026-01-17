package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ibp-network/ibp-geodns-agent/src/agent"
	"github.com/ibp-network/ibp-geodns-agent/src/config"
	"github.com/ibp-network/ibp-geodns-agent/src/logging"
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
	goVersion = runtime.Version()
)

// GetVersion returns version information
func GetVersion() map[string]string {
	return map[string]string{
		"version":   version,
		"buildTime": buildTime,
		"gitCommit": gitCommit,
		"goVersion": goVersion,
	}
}

func main() {
	var (
		configPath  = flag.String("config", "/etc/ibpdns/agent.json", "Path to configuration file")
		showVersion = flag.Bool("version", false, "Show version information")
		logLevel    = flag.String("log-level", "", "Override log level (Debug, Info, Warn, Error)")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("ibp-geodns-agent version %s (built %s, commit %s, go %s)\n", version, buildTime, gitCommit, goVersion)
		os.Exit(0)
	}

	// Initialize logging
	logging.Init(*logLevel)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logging.Fatal("Failed to load configuration", "error", err)
	}

	// Override log level if specified via flag
	if *logLevel != "" {
		cfg.System.LogLevel = *logLevel
	}

	// Set log level from config
	logging.SetLevel(cfg.System.LogLevel)

	logging.Info("Starting ibp-geodns-agent", "version", version, "config", *configPath)

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize and start agent
	a, err := agent.New(cfg)
	if err != nil {
		logging.Fatal("Failed to initialize agent", "error", err)
	}

	// Start agent
	if err := a.Start(ctx); err != nil {
		logging.Fatal("Failed to start agent", "error", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	logging.Info("Received shutdown signal", "signal", sig.String())

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Stop agent gracefully
	if err := a.Stop(shutdownCtx); err != nil {
		logging.Error("Error during agent shutdown", "error", err)
	} else {
		logging.Info("Agent stopped gracefully")
	}
}
