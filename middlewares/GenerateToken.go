package middlewares

import (
	"github.com/google/uuid"
)

func GenerateToken() string {
	token := uuid.New()
	return token.String()
}

