package nats

import (
	"time"

	"github.com/ibp-network/ibp-geodns-agent/src/config"
	"github.com/ibp-network/ibp-geodns-agent/src/logging"
	libsnats "github.com/ibp-network/ibp-geodns-libs/nats"
	natsgo "github.com/nats-io/nats.go"
)

// Client wraps the NATS connection from ibp-geodns-libs
// The libs package uses a global connection pattern, so this wrapper
// provides a cleaner interface for the agent
type Client struct {
	cfg config.NatsConfig
}

// NewClient creates a new NATS client using ibp-geodns-libs
// Note: The libs package uses a global connection, so Connect() must be called
// before using other functions. This wrapper manages that connection.
func NewClient(cfg config.NatsConfig) (*Client, error) {
	// The ibp-geodns-libs/nats package uses a global connection pattern
	// We need to ensure Connect() is called, which uses the global NC variable
	// The connection is configured through the libs package's internal config
	
	// For now, we'll use the libs Connect() function
	// The actual connection details should be set up in the libs package
	if err := libsnats.Connect(); err != nil {
		return nil, err
	}

	logging.Info("Connected to NATS via ibp-geodns-libs", "nodeID", cfg.NodeID)

	return &Client{
		cfg: cfg,
	}, nil
}

// Publish publishes a message to a subject using ibp-geodns-libs
func (c *Client) Publish(subject string, data []byte) error {
	return libsnats.Publish(subject, data)
}

// Subscribe subscribes to a subject using ibp-geodns-libs
func (c *Client) Subscribe(subject string, handler natsgo.MsgHandler) (*natsgo.Subscription, error) {
	// Convert natsgo.MsgHandler to the libs callback format
	cb := func(msg *natsgo.Msg) {
		handler(msg)
	}
	return libsnats.Subscribe(subject, cb)
}

// QueueSubscribe subscribes to a subject with queue group
// Note: The libs package may not have QueueSubscribe, so we use the underlying connection
func (c *Client) QueueSubscribe(subject, queue string, handler natsgo.MsgHandler) (*natsgo.Subscription, error) {
	conn := libsnats.GetConnection()
	if conn == nil {
		return nil, libsnats.Connect()
	}
	cb := func(msg *natsgo.Msg) {
		handler(msg)
	}
	return conn.QueueSubscribe(subject, queue, cb)
}

// Close closes the NATS connection
func (c *Client) Close() error {
	libsnats.Disconnect()
	return nil
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	conn := libsnats.GetConnection()
	return conn != nil && conn.IsConnected()
}

// GetConnection returns the underlying NATS connection for advanced usage
func (c *Client) GetConnection() *natsgo.Conn {
	return libsnats.GetConnection()
}

// Request sends a request and waits for a response
func (c *Client) Request(subject string, data []byte, timeout time.Duration) (*natsgo.Msg, error) {
	return libsnats.Request(subject, data, timeout)
}
