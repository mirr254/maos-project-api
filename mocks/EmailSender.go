package mocks

import (
	"github.com/stretchr/testify/mock"

)

//implement the mock
type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(smtpHost, smtpPort, from, pass, toEmail, subject, body string) error {
	args := m.Called(smtpHost, smtpPort, from, pass, toEmail, subject, body)
	return args.Error(0)
}