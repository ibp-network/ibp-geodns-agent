package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/ibp-network/ibp-geodns-agent/src/config"
	"github.com/ibp-network/ibp-geodns-agent/src/health"
	"github.com/ibp-network/ibp-geodns-agent/src/logging"
	"github.com/ibp-network/ibp-geodns-agent/src/nats"
	"github.com/ibp-network/ibp-geodns-agent/src/reporter"
)

// Agent represents the main agent instance
type Agent struct {
	config   *config.Config
	reporter *reporter.Reporter
	health   *health.Server
	ctx      context.Context
	cancel   context.CancelFunc
}

// New creates a new agent instance
func New(cfg *config.Config) (*Agent, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize NATS connection using ibp-geodns-libs
	if err := nats.Init(cfg.Nats); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to initialize NATS: %w", err)
	}

	// Initialize reporter
	rep, err := reporter.New(cfg)
	if err != nil {
		nats.Disconnect()
		cancel()
		return nil, fmt.Errorf("failed to create reporter: %w", err)
	}

	// Initialize health server
	healthServer := health.New(cfg.Agent.HealthCheckPort)

	return &Agent{
		config:   cfg,
		reporter: rep,
		health:   healthServer,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start starts the agent
func (a *Agent) Start(ctx context.Context) error {
	logging.Info("Starting agent", "agentID", a.config.Agent.AgentID)

	// Start health server
	if err := a.health.Start(); err != nil {
		return fmt.Errorf("failed to start health server: %w", err)
	}

	// Start reporter
	if err := a.reporter.Start(ctx); err != nil {
		return fmt.Errorf("failed to start reporter: %w", err)
	}

	// Start monitoring loop
	go a.monitorLoop(ctx)

	// Start config reload watcher
	go a.configReloadLoop(ctx)

	logging.Info("Agent started successfully")
	return nil
}

// Stop stops the agent gracefully
func (a *Agent) Stop(ctx context.Context) error {
	logging.Info("Stopping agent")

	// Cancel context
	a.cancel()

	// Stop reporter
	if err := a.reporter.Stop(ctx); err != nil {
		logging.Error("Error stopping reporter", "error", err)
	}

	// Stop health server
	if err := a.health.Stop(ctx); err != nil {
		logging.Error("Error stopping health server", "error", err)
	}

	// Close NATS connection
	nats.Disconnect()

	logging.Info("Agent stopped")
	return nil
}

// monitorLoop runs the main monitoring loop
func (a *Agent) monitorLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(a.config.Agent.CheckInterval) * time.Second)
	defer ticker.Stop()

	// Run initial check
	a.performChecks()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			a.performChecks()
		}
	}
}

// performChecks performs health checks on all configured services
func (a *Agent) performChecks() {
	logging.Debug("Performing service checks")

	for _, service := range a.config.Agent.ServicesToMonitor {
		go a.checkService(service)
	}
}

// checkService checks a single service
func (a *Agent) checkService(service config.ServiceConfig) {
	// This would implement actual service checking logic
	// For now, it's a placeholder
	logging.Debug("Checking service", "service", service.Name, "type", service.Type)
	
	// TODO: Implement actual health checks based on service type
	// - HTTP: Make HTTP request, check status code
	// - TCP: Check TCP connection
	// - Custom: Run custom check command
}

// configReloadLoop periodically reloads configuration
func (a *Agent) configReloadLoop(ctx context.Context) {
	if a.config.System.ConfigReloadTime <= 0 {
		return
	}

	ticker := time.NewTicker(time.Duration(a.config.System.ConfigReloadTime) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logging.Debug("Reloading configuration")
			// TODO: Implement config reload logic
		}
	}
}
