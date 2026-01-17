package reporter

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ibp-network/ibp-geodns-agent/src/config"
	"github.com/ibp-network/ibp-geodns-agent/src/logging"
	"github.com/ibp-network/ibp-geodns-agent/src/nats"
)

// Reporter handles reporting agent status and metrics
type Reporter struct {
	config *config.Config
	nats   *nats.Client
	ctx    context.Context
	cancel context.CancelFunc
}

// Report represents a status report
type Report struct {
	AgentID    string                 `json:"agent_id"`
	Timestamp  time.Time              `json:"timestamp"`
	Status     string                 `json:"status"` // online, offline, degraded
	Services   map[string]ServiceStatus `json:"services"`
	Metrics    map[string]interface{} `json:"metrics,omitempty"`
}

// ServiceStatus represents the status of a monitored service
type ServiceStatus struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"` // up, down, degraded
	Latency   time.Duration `json:"latency,omitempty"`
	LastCheck time.Time     `json:"last_check"`
	Error     string        `json:"error,omitempty"`
}

// New creates a new reporter
func New(cfg *config.Config, nc *nats.Client) (*Reporter, error) {
	ctx, cancel := context.WithCancel(context.Background())

	return &Reporter{
		config: cfg,
		nats:   nc,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

// Start starts the reporter
func (r *Reporter) Start(ctx context.Context) error {
	logging.Info("Starting reporter")

	// Start reporting loop
	go r.reportLoop(ctx)

	return nil
}

// Stop stops the reporter
func (r *Reporter) Stop(ctx context.Context) error {
	logging.Info("Stopping reporter")
	r.cancel()
	return nil
}

// reportLoop periodically sends reports
func (r *Reporter) reportLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(r.config.Agent.ReportInterval) * time.Second)
	defer ticker.Stop()

	// Send initial report
	r.sendReport()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.sendReport()
		}
	}
}

// sendReport sends a status report
func (r *Reporter) sendReport() {
	report := Report{
		AgentID:   r.config.Agent.AgentID,
		Timestamp: time.Now(),
		Status:    "online",
		Services:  make(map[string]ServiceStatus),
		Metrics:   make(map[string]interface{}),
	}

	// TODO: Populate service statuses from actual checks
	// For now, this is a placeholder

	data, err := json.Marshal(report)
	if err != nil {
		logging.Error("Failed to marshal report", "error", err)
		return
	}

	// Publish to NATS subject
	subject := fmt.Sprintf("agent.report.%s", r.config.Agent.AgentID)
	if err := r.nats.Publish(subject, data); err != nil {
		logging.Error("Failed to publish report", "error", err, "subject", subject)
		return
	}

	logging.Debug("Sent report", "subject", subject)
}

// ReportServiceStatus reports the status of a specific service
func (r *Reporter) ReportServiceStatus(serviceName string, status ServiceStatus) {
	// This can be called from the monitoring loop to report individual service status
	// For now, it's a placeholder
	logging.Debug("Service status update", "service", serviceName, "status", status.Status)
}
