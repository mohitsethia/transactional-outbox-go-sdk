package outbox

import "time"

// Message represents the message structure in the outbox table
type Message struct {
	ID          uint      `gorm:"primaryKey"`
	Payload     string    `gorm:"type:text"`
	Status      string    `gorm:"type:varchar(50);default:'pending'"`
	ProcessedAt time.Time `gorm:"default:null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
