package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// Config holds the configuration for the NATS connection
type Config struct {
	// Optional existing NATS connection. If provided, we will use it directly.
	NATSConnection *nats.Conn
	// URL for NATS connection if a new connection needs to be created
	URL string
}

// Validate validates the provided NATS configuration
func (c *Config) Validate() error {
	if c.NATSConnection == nil && c.URL == "" {
		return fmt.Errorf("either NATSConnection or URL must be provided")
	}
	return nil
}
