package mock

import (
	"github.com/stretchr/testify/mock"
)

// PublisherMock mocks the Publisher interface
type PublisherMock struct {
	mock.Mock
}

func (m *PublisherMock) PublishMessage(subject string, payload []byte) error {
	args := m.Called(subject, payload)
	return args.Error(0)
}

func (m *PublisherMock) Close() {
	m.Called()
}
