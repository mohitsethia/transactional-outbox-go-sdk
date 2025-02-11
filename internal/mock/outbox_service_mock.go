package mock

import "github.com/stretchr/testify/mock"

// OutboxServiceMock Mocking the Service layer
type OutboxServiceMock struct {
	mock.Mock
}

func (m *OutboxServiceMock) CreateOutboxMessage(payload string) error {
	args := m.Called(payload)
	return args.Error(0)
}

func (m *OutboxServiceMock) ProcessOutboxMessages() error {
	args := m.Called()
	return args.Error(0)
}
