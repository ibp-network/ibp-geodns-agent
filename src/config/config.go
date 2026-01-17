package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/ibp-network/ibp-geodns-libs/config"
)

// Config represents the agent configuration structure
type Config struct {
	System      SystemConfig      `json:"System"`
	Nats        NatsConfig        `json:"Nats"`
	Mysql       MysqlConfig       `json:"Mysql,omitempty"`
	Matrix      MatrixConfig      `json:"Matrix,omitempty"`
	CollatorApi CollatorApiConfig `json:"CollatorApi,omitempty"`
	Agent       AgentConfig       `json:"Agent"`

	mu sync.RWMutex
}

// SystemConfig contains system-level configuration
type SystemConfig struct {
	WorkDir           string                `json:"WorkDir"`
	LogLevel          string                `json:"LogLevel"`
	ConfigUrls        ConfigUrls            `json:"ConfigUrls"`
	ConfigReloadTime  int                   `json:"ConfigReloadTime"`  // seconds
	MinimumOfflineTime int                  `json:"MinimumOfflineTime"` // seconds
}

// ConfigUrls contains URLs for remote configuration
type ConfigUrls struct {
	StaticDNSConfig      string `json:"StaticDNSConfig"`
	MembersConfig        string `json:"MembersConfig"`
	ServicesConfig       string `json:"ServicesConfig"`
	IaasPricingConfig    string `json:"IaasPricingConfig,omitempty"`
	ServicesRequestsConfig string `json:"ServicesRequestsConfig,omitempty"`
}

// NatsConfig contains NATS connection configuration
type NatsConfig struct {
	NodeID string `json:"NodeID"`
	Url    string `json:"Url"`
	User   string `json:"User"`
	Pass   string `json:"Pass"`
}

// MysqlConfig contains MySQL database configuration
type MysqlConfig struct {
	Host string `json:"Host"`
	Port string `json:"Port"`
	User string `json:"User"`
	Pass string `json:"Pass"`
	DB   string `json:"DB"`
}

// MatrixConfig contains Matrix notification configuration
type MatrixConfig struct {
	HomeServerURL string `json:"HomeServerURL"`
	Username      string `json:"Username"`
	Password      string `json:"Password"`
	RoomID        string `json:"RoomID"`
}

// CollatorApiConfig contains collator API configuration
type CollatorApiConfig struct {
	ListenAddress string `json:"ListenAddress"`
	ListenPort    string `json:"ListenPort"`
}

// AgentConfig contains agent-specific configuration
type AgentConfig struct {
	AgentID           string          `json:"AgentID"`
	ReportInterval    int             `json:"ReportInterval"`    // seconds
	CheckInterval     int             `json:"CheckInterval"`     // seconds
	HealthCheckPort   int             `json:"HealthCheckPort"`
	ServicesToMonitor []ServiceConfig `json:"ServicesToMonitor"`
}

// ServiceConfig defines a service to monitor
type ServiceConfig struct {
	Name             string `json:"Name"`
	Type             string `json:"Type"` // http, tcp, custom
	URL              string `json:"URL,omitempty"`
	Endpoint         string `json:"Endpoint,omitempty"`
	Timeout          int    `json:"Timeout"`          // seconds
	Interval         int    `json:"Interval"`         // seconds
	ExpectedStatus   int    `json:"ExpectedStatus,omitempty"`
	ExpectedResponse string `json:"ExpectedResponse,omitempty"`
}

var (
	globalConfig *Config
	configMu     sync.RWMutex
)

// Load loads configuration from file
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Set defaults
	cfg.setDefaults()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Load remote config if URLs are provided
	if err := cfg.loadRemoteConfig(); err != nil {
		// Log warning but don't fail - remote config is optional
		fmt.Printf("Warning: failed to load remote config: %v\n", err)
	}

	configMu.Lock()
	globalConfig = &cfg
	configMu.Unlock()

	return &cfg, nil
}

// Get returns the global configuration (thread-safe)
func Get() *Config {
	configMu.RLock()
	defer configMu.RUnlock()
	return globalConfig
}

// setDefaults sets default values for configuration
func (c *Config) setDefaults() {
	if c.System.WorkDir == "" {
		c.System.WorkDir = "/opt/ibp-geodns-agent"
	}
	if c.System.LogLevel == "" {
		c.System.LogLevel = "Info"
	}
	if c.System.ConfigReloadTime == 0 {
		c.System.ConfigReloadTime = 3600
	}
	if c.System.MinimumOfflineTime == 0 {
		c.System.MinimumOfflineTime = 900
	}
	if c.Agent.ReportInterval == 0 {
		c.Agent.ReportInterval = 60
	}
	if c.Agent.CheckInterval == 0 {
		c.Agent.CheckInterval = 30
	}
	if c.Agent.HealthCheckPort == 0 {
		c.Agent.HealthCheckPort = 8080
	}
	if c.Agent.AgentID == "" {
		hostname, _ := os.Hostname()
		c.Agent.AgentID = hostname
	}
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Nats.Url == "" {
		return fmt.Errorf("NATS URL is required")
	}
	if c.Nats.NodeID == "" {
		return fmt.Errorf("NATS NodeID is required")
	}
	return nil
}

// loadRemoteConfig loads configuration from remote URLs using ibp-geodns-libs
func (c *Config) loadRemoteConfig() error {
	// Use ibp-geodns-libs config loader for remote config
	// For now, this is a placeholder that would integrate with ibp-geodns-libs
	// The actual implementation would use the libs to fetch and merge remote configs
	_ = config.Config{}

	return nil
}

// Reload reloads configuration from file
func (c *Config) Reload(configPath string) error {
	newCfg, err := Load(configPath)
	if err != nil {
		return err
	}

	configMu.Lock()
	*globalConfig = *newCfg
	configMu.Unlock()

	return nil
}
