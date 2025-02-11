package main

import (
	"log"
	"time"

	db "github.com/outbox-go-sdk/internal/db/postgres"
	"github.com/outbox-go-sdk/internal/outbox/handler"
	"github.com/outbox-go-sdk/internal/outbox/service"
	"github.com/outbox-go-sdk/internal/publisher/nats"
)

func main() {
	// Initialize DB config (use existing DB instance or provide connection params)
	dbConfig := &db.Config{
		User:     "postgres",     // Use the user from docker-compose.yml
		Password: "rootpassword", // Use the password from docker-compose.yml
		Host:     "postgres",     // Use the service name from docker-compose.yml
		Port:     5432,           // Default PostgreSQL port
		DBName:   "transactional_outbox",
		SSLMode:  "disable", // Optional: default is "disable"
	}

	// Initialize NATS config (use existing NATS connection or provide URL)
	natsConfig := &nats.Config{
		URL: "nats://nats:4222", // Use the NATS URL
	}

	// Initialize Repositories with Config structs
	dbRepo, err := db.NewGormRepository(dbConfig)
	if err != nil {
		log.Fatalf("Error initializing DB: %v", err)
	}

	ncRepo, err := nats.NewNatsPublisher(natsConfig)
	if err != nil {
		log.Fatalf("Error initializing NATS: %v", err)
	}

	// Initialize Service with the repositories
	outboxService := service.NewService(dbRepo, ncRepo, 100)

	// Initialize Handler with the service
	outboxHandler := handler.NewHandler(outboxService)

	// Start processing outbox messages
	for {
		outboxHandler.Process()

		// Sleep before checking again
		time.Sleep(2 * time.Second)
	}
}
