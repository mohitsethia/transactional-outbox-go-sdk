package postgres

import (
	"time"

	"github.com/outbox-go-sdk/internal/domain/outbox"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository interface {
	// Methods to interact with the database
	CreateOutboxMessage(message outbox.Message) error
	FindUnprocessedMessages(batchSize int) ([]outbox.Message, error)
	MarkMessageAsProcessed(message outbox.Message) error
	BeginTransaction() Repository
	RollBackTransaction() error
	CommitTransaction() error
}

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(config *Config) (Repository, error) {
	var db *gorm.DB
	var err error

	// Use provided DB instance or create a new connection
	if config.DBInstance != nil {
		db = config.DBInstance
	} else {
		dsn := config.BuildDSN()
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
	}

	// Auto-migrate the OutboxMessage model
	if err := db.AutoMigrate(&outbox.Message{}); err != nil {
		return nil, err
	}

	return &gormRepository{db: db}, nil
}

// CreateOutboxMessage adds a new message to the outbox table
func (r *gormRepository) CreateOutboxMessage(message outbox.Message) error {
	if err := r.db.Create(&message).Error; err != nil {
		return err
	}
	return nil
}

// BeginTransaction starts a new database transaction
func (r *gormRepository) BeginTransaction() Repository {
	return &gormRepository{
		db: r.db.Begin(),
	}
}

// FindUnprocessedMessages retrieves unprocessed outbox messages in batches
func (r *gormRepository) FindUnprocessedMessages(batchSize int) ([]outbox.Message, error) {
	var messages []outbox.Message
	if err := r.db.Where("status = ?", "pending").Limit(batchSize).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// MarkMessageAsProcessed marks a message as processed in the database
func (r *gormRepository) MarkMessageAsProcessed(message outbox.Message) error {
	processedAt := time.Now()
	if err := r.db.Model(&message).UpdateColumns(map[string]interface{}{
		"status":       "processed",
		"processed_at": processedAt,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *gormRepository) CommitTransaction() error {
	return r.db.Commit().Error
}

func (r *gormRepository) RollBackTransaction() error {
	return r.db.Rollback().Error
}
