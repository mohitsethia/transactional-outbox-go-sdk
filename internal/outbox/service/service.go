package service

import (
	"log"

	db "github.com/outbox-go-sdk/internal/db/postgres"
	"github.com/outbox-go-sdk/internal/domain/outbox"
	"github.com/outbox-go-sdk/internal/publisher/nats"
)

type Service interface {
	CreateOutboxMessage(payload string) error
	ProcessOutboxMessages() error
}

// Service defines the logic of handling outbox messages
type service struct {
	dbRepo    db.Repository
	msgRepo   nats.Publisher
	batchSize int
}

// NewService creates a new instance of Service
func NewService(dbRepo db.Repository, msgRepo nats.Publisher, batchSize int) Service {
	return &service{
		dbRepo:    dbRepo,
		msgRepo:   msgRepo,
		batchSize: batchSize,
	}
}

// CreateOutboxMessage creates a new message and adds it to the outbox table
func (s *service) CreateOutboxMessage(payload string) error {
	// Create the outbox message
	message := outbox.Message{
		Payload: payload,
	}

	dbRepo := s.dbRepo.BeginTransaction()
	var err error
	defer func() {
		if err != nil {
			err = dbRepo.RollBackTransaction()
			if err != nil {
				log.Printf("Error rolling back transaction: %v", err)
			}
		}
	}()

	// Add the message to the database (outbox table)
	if err = dbRepo.CreateOutboxMessage(message); err != nil {
		log.Printf("Error creating outbox message: %v", err)
		return err
	}

	// Commit the transaction
	if err = dbRepo.CommitTransaction(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}

// ProcessOutboxMessages retrieves unprocessed messages, publishes them, and marks them as processed
func (s *service) ProcessOutboxMessages() error {
	// Start a database transaction
	dbRepo := s.dbRepo.BeginTransaction()

	var err error
	defer func() {
		if err != nil {
			err = dbRepo.RollBackTransaction()
			if err != nil {
				log.Printf("Error rolling back outbox message: %v", err)
			}
		}
	}()

	// Retrieve unprocessed outbox messages
	messages, err := dbRepo.FindUnprocessedMessages(s.batchSize)
	if err != nil {
		log.Printf("Error fetching unprocessed messages: %v", err)
		return err
	}

	// Process each message within the transaction
	for _, message := range messages {
		// Publish to NATS
		if err = s.msgRepo.PublishMessage("outbox", []byte(message.Payload)); err != nil {
			log.Printf("Error publishing message: %v", err)
			return err
		}

		// Mark the message as processed
		if err = dbRepo.MarkMessageAsProcessed(message); err != nil {
			log.Printf("Error marking message as processed: %v", err)
			return err
		}
	}

	// Commit the transaction
	if err = dbRepo.CommitTransaction(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		return err
	}

	return nil
}
