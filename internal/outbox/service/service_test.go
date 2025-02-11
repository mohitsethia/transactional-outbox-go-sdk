package service

import (
	"errors"
	"testing"

	"github.com/outbox-go-sdk/internal/domain/outbox"
	mock2 "github.com/outbox-go-sdk/internal/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateOutboxMessage_Success(t *testing.T) {
	mockDB := new(mock2.DBRepoMock)
	mockPublisher := new(mock2.PublisherMock)
	service := NewService(mockDB, mockPublisher, 10)

	mockDB.On("BeginTransaction").Return(mockDB)
	mockDB.On("CreateOutboxMessage", mock.Anything).Return(nil)
	mockDB.On("CommitTransaction").Return(nil)

	err := service.CreateOutboxMessage("Test Payload")
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
}

func TestCreateOutboxMessage_Failure_CreateError(t *testing.T) {
	mockDB := new(mock2.DBRepoMock)
	mockPublisher := new(mock2.PublisherMock)
	service := NewService(mockDB, mockPublisher, 10)

	mockDB.On("BeginTransaction").Return(mockDB)
	mockDB.On("CreateOutboxMessage", mock.Anything).Return(errors.New("db error"))
	mockDB.On("RollBackTransaction").Return(nil)

	err := service.CreateOutboxMessage("Test Payload")
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
}

func TestProcessOutboxMessages_Success(t *testing.T) {
	mockDB := new(mock2.DBRepoMock)
	mockPublisher := new(mock2.PublisherMock)
	service := NewService(mockDB, mockPublisher, 10)

	mockDB.On("BeginTransaction").Return(mockDB)
	mockDB.On("FindUnprocessedMessages", 10).Return([]outbox.Message{
		{Payload: "Test Payload"},
	}, nil)
	mockPublisher.On("PublishMessage", "outbox", mock.Anything).Return(nil)
	mockDB.On("MarkMessageAsProcessed", mock.Anything).Return(nil)
	mockDB.On("CommitTransaction").Return(nil)

	err := service.ProcessOutboxMessages()
	assert.NoError(t, err)
	mockDB.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestProcessOutboxMessages_Failure_FetchError(t *testing.T) {
	mockDB := new(mock2.DBRepoMock)
	mockPublisher := new(mock2.PublisherMock)
	service := NewService(mockDB, mockPublisher, 10)

	mockDB.On("BeginTransaction").Return(mockDB)
	mockDB.On("FindUnprocessedMessages", 10).Return(([]outbox.Message)(nil), errors.New("db error"))
	mockDB.On("RollBackTransaction").Return(nil)

	err := service.ProcessOutboxMessages()
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}

func TestProcessOutboxMessages_Failure_PublishError(t *testing.T) {
	mockDB := new(mock2.DBRepoMock)
	mockPublisher := new(mock2.PublisherMock)
	service := NewService(mockDB, mockPublisher, 10)

	mockDB.On("BeginTransaction").Return(mockDB)
	mockDB.On("FindUnprocessedMessages", 10).Return([]outbox.Message{
		{Payload: "Test Payload"},
	}, nil)
	mockPublisher.On("PublishMessage", "outbox", mock.Anything).Return(errors.New("nats error"))
	mockDB.On("RollBackTransaction").Return(nil)

	err := service.ProcessOutboxMessages()
	assert.Error(t, err)
	mockDB.AssertExpectations(t)
	mockPublisher.AssertExpectations(t)
}
