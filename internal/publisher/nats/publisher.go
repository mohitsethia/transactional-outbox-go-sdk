package nats

import (
	"github.com/nats-io/nats.go"
)

// Publisher defines methods for interacting with NATS
type Publisher interface {
	PublishMessage(subject string, data []byte) error
	Close()
}

// publisher implements the Publisher interface using NATS
type publisher struct {
	nc *nats.Conn
}

// NewNatsPublisher creates a new instance of NatsPublisher using the provided config
func NewNatsPublisher(config *Config) (Publisher, error) {
	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	var nc *nats.Conn
	var err error

	// Use the provided NATS connection if available
	if config.NATSConnection != nil {
		nc = config.NATSConnection
	} else {
		// Create a new NATS connection using the provided URL
		nc, err = nats.Connect(config.URL)
		if err != nil {
			return nil, err
		}
	}

	return &publisher{nc: nc}, nil
}

// PublishMessage sends a message to a NATS subject
func (r *publisher) PublishMessage(subject string, data []byte) error {
	if err := r.nc.Publish(subject, data); err != nil {
		return err
	}
	return nil
}

func (r *publisher) Close() {
	r.nc.Close()
}
