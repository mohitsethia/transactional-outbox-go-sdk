package handler

import (
	"log"

	"github.com/outbox-go-sdk/internal/outbox/service"
)

type Handler interface {
	CreateMessage(payload string) error
	Process()
}

// Handler is the entry point to initiate processing of outbox messages
type handler struct {
	service service.Service
}

// NewHandler initializes a new Handler
func NewHandler(service service.Service) Handler {
	return &handler{service: service}
}

func (h *handler) CreateMessage(payload string) error {
	if err := h.service.CreateOutboxMessage(payload); err != nil {
		log.Printf("Error creating outbox message: %v", err)
		return err
	}
	return nil
}

// Process initiates the outbox message processing
func (h *handler) Process() {
	// Call the internal service method to process the messages
	if err := h.service.ProcessOutboxMessages(); err != nil {
		log.Printf("Error processing outbox messages: %v", err)
	}
}
