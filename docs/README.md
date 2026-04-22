# ibp-geodns-agent

A self-hosted active agent for monitoring ancillary services and reporting SLA metrics to the IBP GeoDNS monitoring infrastructure.

## Overview

`ibp-geodns-agent` is an early-stage system service that runs on machines to be monitored. The current build provides NATS connectivity, periodic self-report publication, and health endpoints; the service-checking and remote-config portions are scaffolding and should not be treated as production-ready features yet.

## Features

- **Health Endpoints**: Provides HTTP health check endpoints for orchestration
- **NATS Connectivity**: Maintains a long-lived NATS connection for future agent integrations
- **Periodic Self-Reporting**: Publishes heartbeat-style agent reports on a configurable interval
- **Structured Logging**: Key/value logging with configurable levels
- **Scaffolded Monitoring**: Service monitoring hooks exist, but concrete HTTP/TCP/custom checks are not implemented in this build
- **Scaffolded Remote Config**: Remote configuration fields are parsed, but remote config fetching/merging is not implemented yet
- **System Service**: Installable as a systemd service

## Requirements

- Go 1.21 or later
- NATS server access
- Linux system (for systemd service)

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/ibp-network/ibp-geodns-agent.git
cd ibp-geodns-agent

# Build
make build

# Install system-wide
sudo make install
```

### System Service Installation

```bash
# Copy systemd service file
sudo cp docs/systemd/ibpagent.systemd /etc/systemd/system/ibp-agent.service

# Create configuration directory
sudo mkdir -p /etc/ibpdns

# Copy and edit configuration
sudo cp config/config.json /etc/ibpdns/agent.json
sudo nano /etc/ibpdns/agent.json

# Enable and start service
sudo systemctl daemon-reload
sudo systemctl enable ibp-agent
sudo systemctl start ibp-agent

# Check status
sudo systemctl status ibp-agent
```

## Configuration

Configuration is provided via a JSON file. See `config/config.json` for an example configuration.

### Configuration Structure

```json
{
  "System": {
    "WorkDir": "/opt/ibp-geodns-agent",
    "LogLevel": "Info",
    "ConfigUrls": {
      "StaticDNSConfig": "https://...",
      "MembersConfig": "https://...",
      "ServicesConfig": "https://..."
    },
    "ConfigReloadTime": 3600,
    "MinimumOfflineTime": 900
  },
  "Nats": {
    "NodeID": "agent-1",
    "Url": "nats://127.0.0.1:4222",
    "User": "natsuser",
    "Pass": "natspasswd"
  },
  "Agent": {
    "AgentID": "agent-1",
    "ReportInterval": 60,
    "CheckInterval": 30,
    "HealthCheckPort": 8080,
    "ServicesToMonitor": [
      {
        "Name": "RPC Service",
        "Type": "http",
        "URL": "https://rpc.example.com/health",
        "Timeout": 10,
        "Interval": 60,
        "ExpectedStatus": 200
      }
    ]
  }
}
```

### Configuration Options

- **System.WorkDir**: Working directory for the agent
- **System.LogLevel**: Logging level (Debug, Info, Warn, Error, Fatal)
- **System.ConfigReloadTime**: Interval in seconds reserved for future remote configuration reload support
- **Nats**: NATS connection configuration
- **Agent.AgentID**: Unique identifier for this agent instance
- **Agent.ReportInterval**: Interval in seconds between status reports
- **Agent.CheckInterval**: Interval in seconds reserved for future service checks
- **Agent.ServicesToMonitor**: Service definitions for the planned monitoring subsystem

## Usage

### Command Line Options

```bash
ibp-agent --config /path/to/config.json
ibp-agent --version
ibp-agent --log-level Debug
```

### Health Endpoints

The agent exposes HTTP health check endpoints:

- `GET /health` - Health check (returns 200 if healthy)
- `GET /ready` - Readiness check (returns 200 if ready)
- `GET /live` - Liveness check (always returns 200 if running)

Default port is 8080, configurable via `Agent.HealthCheckPort`.

## Architecture

The agent follows the same structural patterns as other IBP GeoDNS repositories:

- **src/main.go**: Main entry point and CLI
- **src/config/**: Configuration loading and management
- **src/agent/**: Core agent logic
- **src/nats/**: NATS client wrapper using ibp-geodns-libs
- **src/reporter/**: Periodic self-report publishing
- **src/health/**: Health check server
- **src/logging/**: Structured logging

## Integration with ibp-geodns-libs

This agent uses `ibp-geodns-libs` for:

- NATS connectivity and messaging
- Shared config structures used by the wider IBP GeoDNS stack
- Common utilities and patterns

## Development

```bash
# Run tests
make test

# Format code
make fmt

# Run linters
make lint

# Build
make build

# Run locally
make run
```

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please ensure your code follows the same patterns and conventions used in other IBP GeoDNS repositories.
