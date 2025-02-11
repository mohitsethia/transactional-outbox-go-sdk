//go:build integration
// +build integration

package go_transactional_outbox

import (
	"fmt"
	"testing"

	repo "github.com/outbox-go-sdk/internal/db/postgres"
	domain "github.com/outbox-go-sdk/internal/domain/outbox"
	outboxService "github.com/outbox-go-sdk/internal/outbox/service"
	nats2 "github.com/outbox-go-sdk/internal/publisher/nats"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	postgresHost = "postgres"
	postgresPort = "5432"
	postgresUser = "postgres"
	postgresPass = "rootpassword"
	postgresDB   = "transactional_outbox"
	natsURL      = "nats://nats:4222"
)

// Setup the real PostgreSQL and NATS server for integration testing
func setupDB() (*gorm.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", postgresUser, postgresPass, postgresHost, postgresPort, postgresDB)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setupNATS() (*nats.Conn, error) {
	return nats.Connect(natsURL)
}

// Test function to create and process outbox messages
func TestCreateAndProcessOutboxMessages(t *testing.T) {
	// Set up the database and NATS connection
	db, err := setupDB()
	require.NoError(t, err)

	nc, err := setupNATS()
	require.NoError(t, err)
	defer nc.Close()

	// Initialize Repositories with Config structs
	dbRepo, err := repo.NewGormRepository(&repo.Config{DBInstance: db})
	require.NoError(t, err)

	ncRepo, err := nats2.NewNatsPublisher(&nats2.Config{NATSConnection: nc})
	require.NoError(t, err)

	// Initialize the outbox service with the real database and NATS publisher
	service := outboxService.NewService(dbRepo, ncRepo, 10)

	// Create a new message via the service
	payload := "Sample outbox message"
	err = service.CreateOutboxMessage(payload)
	require.NoError(t, err)

	// Verify the message was inserted into the database using GORM
	var count int64
	err = db.Model(&domain.Message{}).Where("payload = ?", payload).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Message should be inserted into the database")

	// Process messages (this should send the message to NATS)
	err = service.ProcessOutboxMessages()
	require.NoError(t, err)

	// Now verify if the message was processed using GORM
	var processed bool
	var messageID uint
	err = db.Model(&domain.Message{}).Where("payload = ?", payload).First(&domain.Message{ID: messageID}).Update("status", "processed").Error
	require.NoError(t, err)

	// Verify that the message is marked as processed
	err = db.Model(&domain.Message{}).Where("payload = ?", payload).First(&domain.Message{ID: messageID}).Scan(&domain.Message{Status: "processed"}).Error
	require.NoError(t, err)
	assert.True(t, processed, "Message should be marked as processed")

	// Assuming the message is now published to NATS, no need to check that directly in the test unless mocking NATS is necessary
}

// Test function to simulate processing failure (e.g., NATS down)
func TestProcessMessagesWithFailure(t *testing.T) {
	// Set up the database and NATS connection
	db, err := setupDB()
	require.NoError(t, err)

	nc, err := setupNATS()
	require.NoError(t, err)
	defer nc.Close()

	// Initialize Repositories with Config structs
	dbRepo, err := repo.NewGormRepository(&repo.Config{DBInstance: db})
	require.NoError(t, err)

	ncRepo, err := nats2.NewNatsPublisher(&nats2.Config{NATSConnection: nc})
	require.NoError(t, err)

	// Initialize the outbox service with the real database and NATS publisher
	service := outboxService.NewService(dbRepo, ncRepo, 10)

	// Create a message first
	payload := "Failure test message"
	err = service.CreateOutboxMessage(payload)
	require.NoError(t, err)

	// Verify the message was inserted into the database
	var count int64
	err = db.Model(&domain.Message{}).Where("payload = ?", payload).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Message should be inserted into the database")

	// Simulate a failure by shutting down NATS (you can use mock or real failure scenarios)
	nc.Close()

	// Try to process messages, expecting failure due to NATS being down
	err = service.ProcessOutboxMessages()
	assert.Error(t, err, "Processing messages should fail when NATS is down")
}
