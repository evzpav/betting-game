package domain

import (
	"math/rand"
	"time"

	uuid "github.com/satori/go.uuid"
)

func GenerateRandomNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min //nolint:gosec
}

func GenerateNewID() string {
	sID := uuid.NewV1()
	return sID.String()
}
