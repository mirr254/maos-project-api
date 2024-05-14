package mocks

import (
	"maos-cloud-project-api/config"

	"github.com/stretchr/testify/mock"
)

//implement the mock
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail( cfg *config.Config, toEmail, subject, body string) error {
	args := m.Called( cfg, toEmail, subject, body)
	return args.Error(0)
}
