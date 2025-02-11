package handler

import (
	"errors"
	"testing"

	"github.com/outbox-go-sdk/internal/mock"

	"github.com/stretchr/testify/assert"
)

func TestCreateMessage_Success(t *testing.T) {
	mockService := new(mock.OutboxServiceMock)
	handler := NewHandler(mockService)

	mockService.On("CreateOutboxMessage", "Test Payload").Return(nil)

	err := handler.CreateMessage("Test Payload")
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestCreateMessage_Failure(t *testing.T) {
	mockService := new(mock.OutboxServiceMock)
	handler := NewHandler(mockService)

	mockService.On("CreateOutboxMessage", "Test Payload").Return(errors.New("error creating message"))

	err := handler.CreateMessage("Test Payload")
	assert.Error(t, err)
	mockService.AssertExpectations(t)
}

func TestProcess_Success(t *testing.T) {
	mockService := new(mock.OutboxServiceMock)
	handler := NewHandler(mockService)

	mockService.On("ProcessOutboxMessages").Return(nil)

	handler.Process()
	mockService.AssertExpectations(t)
}

func TestProcess_Failure(t *testing.T) {
	mockService := new(mock.OutboxServiceMock)
	handler := NewHandler(mockService)

	mockService.On("ProcessOutboxMessages").Return(errors.New("error processing messages"))

	handler.Process()
	mockService.AssertExpectations(t)
}
