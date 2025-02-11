package postgres

import (
	"fmt"

	"gorm.io/gorm"
)

// Config holds the configuration for the database connection
type Config struct {
	// Optional existing DB instance. If provided, we will use it directly.
	DBInstance *gorm.DB
	// Database connection parameters for building the DSN string
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
	SSLMode  string // Optional: set it to "disable" if you don't need SSL
}

// Validate validates the provided database configuration
func (c *Config) Validate() error {
	// Validate if DBInstance is provided or if all connection params are provided
	if c.DBInstance == nil && (c.User == "" || c.Password == "" || c.Host == "" || c.Port == 0 || c.DBName == "") {
		return fmt.Errorf("either DBInstance or all database connection params (User, Password, Host, Port, DBName) must be provided")
	}
	return nil
}

// BuildDSN constructs the DSN string from the provided configuration
func (c *Config) BuildDSN() string {
	// Default SSLMode is "disable" if not provided
	if c.SSLMode == "" {
		c.SSLMode = "disable"
	}

	// Build DSN (Data Source Name) string for PostgreSQL
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}
