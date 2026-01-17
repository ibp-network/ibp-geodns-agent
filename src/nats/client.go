package nats

import (
	"fmt"
	"time"

	"github.com/ibp-network/ibp-geodns-agent/src/config"
	"github.com/ibp-network/ibp-geodns-agent/src/logging"
	natsgo "github.com/nats-io/nats.go"
)

// Client wraps the NATS connection
type Client struct {
	conn *natsgo.Conn
	cfg  config.NatsConfig
}

// NewClient creates a new NATS client using ibp-geodns-libs
func NewClient(cfg config.NatsConfig) (*Client, error) {
	// Use ibp-geodns-libs for NATS connectivity
	// For now, this is a placeholder that would use the actual libs implementation
	
	opts := []natsgo.Option{
		natsgo.Name(cfg.NodeID),
		natsgo.ReconnectWait(2 * time.Second),
		natsgo.MaxReconnects(-1),
		natsgo.DisconnectErrHandler(func(nc *natsgo.Conn, err error) {
			if err != nil {
				logging.Warn("NATS disconnected", "error", err)
			}
		}),
		natsgo.ReconnectHandler(func(nc *natsgo.Conn) {
			logging.Info("NATS reconnected", "url", nc.ConnectedUrl())
		}),
	}

	// Add authentication if provided
	if cfg.User != "" && cfg.Pass != "" {
		opts = append(opts, natsgo.UserInfo(cfg.User, cfg.Pass))
	}

	conn, err := natsgo.Connect(cfg.Url, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	logging.Info("Connected to NATS", "url", cfg.Url, "nodeID", cfg.NodeID)

	return &Client{
		conn: conn,
		cfg:  cfg,
	}, nil
}

// Publish publishes a message to a subject
func (c *Client) Publish(subject string, data []byte) error {
	return c.conn.Publish(subject, data)
}

// Subscribe subscribes to a subject
func (c *Client) Subscribe(subject string, handler natsgo.MsgHandler) (*natsgo.Subscription, error) {
	return c.conn.Subscribe(subject, handler)
}

// QueueSubscribe subscribes to a subject with queue group
func (c *Client) QueueSubscribe(subject, queue string, handler natsgo.MsgHandler) (*natsgo.Subscription, error) {
	return c.conn.QueueSubscribe(subject, queue, handler)
}

// Close closes the NATS connection
func (c *Client) Close() error {
	if c.conn != nil {
		c.conn.Close()
	}
	return nil
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}
