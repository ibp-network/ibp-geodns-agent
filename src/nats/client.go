package nats

import (
	"fmt"
	"sync"
	"time"

	"github.com/ibp-network/ibp-geodns-agent/src/config"
	"github.com/ibp-network/ibp-geodns-agent/src/logging"
	natsgo "github.com/nats-io/nats.go"
)

var (
	connMu sync.RWMutex
	conn   *natsgo.Conn
)

func currentConnection() *natsgo.Conn {
	connMu.RLock()
	defer connMu.RUnlock()
	return conn
}

// Init initializes the NATS connection using the agent configuration directly.
func Init(cfg config.NatsConfig) error {
	if cfg.Url == "" {
		return fmt.Errorf("NATS URL is required")
	}

	connMu.Lock()
	defer connMu.Unlock()

	if conn != nil && !conn.IsClosed() {
		return nil
	}

	opts := []natsgo.Option{
		natsgo.Name("ibp-geodns-agent:" + cfg.NodeID),
		natsgo.MaxReconnects(-1),
		natsgo.ReconnectWait(2 * time.Second),
		natsgo.Timeout(10 * time.Second),
		natsgo.DisconnectErrHandler(func(_ *natsgo.Conn, err error) {
			logging.Error("NATS disconnected", "error", err)
		}),
		natsgo.ReconnectHandler(func(c *natsgo.Conn) {
			logging.Info("NATS reconnected", "url", c.ConnectedUrl())
		}),
		natsgo.ClosedHandler(func(c *natsgo.Conn) {
			if err := c.LastError(); err != nil {
				logging.Error("NATS connection closed", "error", err)
			}
		}),
	}
	if cfg.User != "" || cfg.Pass != "" {
		opts = append(opts, natsgo.UserInfo(cfg.User, cfg.Pass))
	}

	connected, err := natsgo.Connect(cfg.Url, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to NATS: %w", err)
	}
	conn = connected

	logging.Info("Connected to NATS", "nodeID", cfg.NodeID, "url", connected.ConnectedUrl())
	return nil
}

// Publish publishes a message to a subject.
func Publish(subject string, data []byte) error {
	active := currentConnection()
	if active == nil || active.IsClosed() {
		return natsgo.ErrConnectionClosed
	}
	return active.Publish(subject, data)
}

// Subscribe subscribes to a subject.
func Subscribe(subject string, cb func(*natsgo.Msg)) (*natsgo.Subscription, error) {
	active := currentConnection()
	if active == nil || active.IsClosed() {
		return nil, natsgo.ErrConnectionClosed
	}
	return active.Subscribe(subject, func(msg *natsgo.Msg) {
		go cb(msg)
	})
}

// Request sends a request and waits for a response.
func Request(subject string, data []byte, timeout time.Duration) (*natsgo.Msg, error) {
	active := currentConnection()
	if active == nil || active.IsClosed() {
		return nil, natsgo.ErrConnectionClosed
	}
	return active.Request(subject, data, timeout)
}

// GetConnection returns the underlying NATS connection for advanced usage.
func GetConnection() *natsgo.Conn {
	return currentConnection()
}

// Disconnect closes the NATS connection.
func Disconnect() {
	connMu.Lock()
	defer connMu.Unlock()
	if conn != nil && !conn.IsClosed() {
		conn.Close()
	}
	conn = nil
}

// IsConnected returns whether the client is connected.
func IsConnected() bool {
	active := currentConnection()
	return active != nil && active.IsConnected()
}
