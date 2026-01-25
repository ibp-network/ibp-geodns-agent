package nats

import (
	"time"

	"github.com/ibp-network/ibp-geodns-agent/src/config"
	"github.com/ibp-network/ibp-geodns-agent/src/logging"
	libsnats "github.com/ibp-network/ibp-geodns-libs/nats"
	natsgo "github.com/nats-io/nats.go"
)

// Init initializes the NATS connection using ibp-geodns-libs
// This should be called once at startup with the agent's NATS configuration
// The libs package uses a global connection, so this sets it up for the agent
func Init(cfg config.NatsConfig) error {
	// The ibp-geodns-libs/nats package uses a global connection
	// The Connect() function reads configuration from the libs config package
	// We need to ensure the libs config is set up with our NATS settings first
	
	// TODO: Configure the libs config package with our NATS settings
	// This might require setting up the libs config with the agent's NATS config
	
	// Connect using the libs package
	if err := libsnats.Connect(); err != nil {
		return err
	}

	logging.Info("Connected to NATS via ibp-geodns-libs", "nodeID", cfg.NodeID)
	return nil
}

// Publish publishes a message to a subject using ibp-geodns-libs
func Publish(subject string, data []byte) error {
	return libsnats.Publish(subject, data)
}

// Subscribe subscribes to a subject using ibp-geodns-libs
func Subscribe(subject string, cb func(*libsnats.NatsMsg)) (*natsgo.Subscription, error) {
	return libsnats.Subscribe(subject, cb)
}

// Request sends a request and waits for a response using ibp-geodns-libs
func Request(subject string, data []byte, timeout time.Duration) (*libsnats.NatsMsg, error) {
	return libsnats.Request(subject, data, timeout)
}

// GetConnection returns the underlying NATS connection for advanced usage
func GetConnection() *natsgo.Conn {
	return libsnats.GetConnection()
}

// Disconnect closes the NATS connection
func Disconnect() {
	libsnats.Disconnect()
}

// IsConnected returns whether the client is connected
func IsConnected() bool {
	conn := libsnats.GetConnection()
	return conn != nil && conn.IsConnected()
}
