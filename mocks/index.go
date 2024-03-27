package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockHarsher struct {
    mock.Mock
}

func (m *MockHarsher) GenerateHashPassword(password string) (string, error) {
    args := m.Called(password)
    return args.String(0), args.Error(1)
}