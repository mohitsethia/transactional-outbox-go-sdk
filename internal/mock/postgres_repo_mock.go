package mock

import (
	"github.com/outbox-go-sdk/internal/db/postgres"
	"github.com/outbox-go-sdk/internal/domain/outbox"

	"github.com/stretchr/testify/mock"
)

// DBRepoMock Mocking db.Repository interface
type DBRepoMock struct {
	mock.Mock
}

func (m *DBRepoMock) BeginTransaction() postgres.Repository {
	args := m.Called()
	return args.Get(0).(postgres.Repository)
}

func (m *DBRepoMock) CreateOutboxMessage(message outbox.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *DBRepoMock) FindUnprocessedMessages(batchSize int) ([]outbox.Message, error) {
	args := m.Called(batchSize)
	return args.Get(0).([]outbox.Message), args.Error(1)
}

func (m *DBRepoMock) MarkMessageAsProcessed(message outbox.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *DBRepoMock) CommitTransaction() error {
	args := m.Called()
	return args.Error(0)
}

func (m *DBRepoMock) RollBackTransaction() error {
	args := m.Called()
	return args.Error(0)
}
