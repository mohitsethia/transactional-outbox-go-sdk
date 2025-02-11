# Transactional Outbox Go SDK

Outbox Go SDK provides a simple mechanism to implement outbox pattern in Go applications, enabling reliable event-driven architectures using PostgreSQL and NATS.

This SDK provides functionality to store events in a transactional outbox and asynchronously process them to send messages to an external message queue (NATS).

## Table of Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Usage](#usage)
- [Docker Setup](#docker-setup)
- [Running Tests](#running-tests)

## Getting Started

To use the `outbox-go-sdk`, you need to have Docker, Docker Compose, and Go installed. The SDK integrates with PostgreSQL as a storage and NATS for message publishing.

### Clone the repository:

```
git clone https://github.com/yourusername/outbox-go-sdk.git
cd outbox-go-sdk
```

### Build the Docker containers:
`make run`\
This command builds and starts the necessary services, including:

`postgres`: PostgreSQL database for storing the outbox messages\
`nats`: NATS server for message publishing\
`app`: Your application with the SDK integrated\

Once the containers are running, your application will automatically connect to the PostgreSQL and NATS services.

## Prerequisites
PostgreSQL: We use PostgreSQL to store events in the outbox table.\
NATS: We use NATS to publish messages once they're processed.\
Docker & Docker Compose: These are used to set up the development environment.\

### Project Structure
The project is structured as follows:
```
.
├── db/                  # Database related files
│   ├── init.sql         # SQL script to initialize tables
├── internal/            # Internal packages
│   ├── db/              # DB-related logic (using GORM)
│   ├── outbox/          # Outbox service logic
│   └── publisher/       # NATS publisher logic
├── Dockerfile           # Dockerfile to build the app container
├── docker-compose.yml   # Docker Compose file for setting up services
├── go.mod               # Go Modules file
├── go.sum               # Go Modules sum file
├── main.go              # Entry point for the application
└── README.md            # This file
```
### Usage
1. Create an Outbox Message
To create a new message and store it in the outbox table, you can use the CreateOutboxMessage method in the outboxService:

```
package main

import (
	"fmt"
	"log"
	"github.com/outbox-go-sdk/internal/outbox/service"
)

func main() {
	// Assuming that dbRepo and ncRepo are properly initialized
	service := outboxService.NewService(dbRepo, ncRepo, 10)

	// Create a new message with some payload
	payload := "Sample outbox message"
	err := service.CreateOutboxMessage(payload)
	if err != nil {
		log.Fatal("Error creating outbox message:", err)
	}
	fmt.Println("Outbox message created successfully!")
}
```
2. Process Outbox Messages
To process the messages in the outbox and publish them to NATS:

```
err := service.ProcessOutboxMessages()
if err != nil {
	log.Fatal("Error processing outbox messages:", err)
}
fmt.Println("Outbox messages processed and published to NATS!")
```

### Docker Setup
The docker-compose.yml file is configured to run the necessary services for PostgreSQL, NATS, and your Go application.

Starting the services:
Run the following command to start the application, database, and NATS:

`make run`

Stopping the services:
To stop the services:

`docker-compose down`

### Running Tests
You can run the unit tests for the SDK by using Go's testing framework.

`make tests`
