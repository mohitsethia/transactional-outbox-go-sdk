package nats

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	mockNats "github.com/outbox-go-sdk/internal/mock" // Ensure correct import path
)

func TestPublishMessage_Success(t *testing.T) {
	mockPublisher := new(mockNats.PublisherMock)

	subject := "test.subject"
	message := []byte("test message")

	mockPublisher.On("PublishMessage", subject, message).Return(nil)

	err := mockPublisher.PublishMessage(subject, message)

	assert.NoError(t, err)
	mockPublisher.AssertExpectations(t)
}

func TestPublishMessage_Failure(t *testing.T) {
	mockPublisher := new(mockNats.PublisherMock)

	subject := "test.subject"
	message := []byte("test message")

	mockPublisher.On("PublishMessage", subject, message).Return(errors.New("publish error"))

	err := mockPublisher.PublishMessage(subject, message)

	assert.Error(t, err)
	assert.Equal(t, "publish error", err.Error())
	mockPublisher.AssertExpectations(t)
}

func TestPublisher_Close(t *testing.T) {
	mockPublisher := new(mockNats.PublisherMock)
	mockPublisher.On("Close").Return()

	mockPublisher.Close()

	mockPublisher.AssertExpectations(t)
}
