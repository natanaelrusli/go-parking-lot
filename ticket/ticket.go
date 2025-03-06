package ticket

import "github.com/google/uuid"

func GenerateTicketNumber() string {
	// Using first 8 characters of UUID
	return uuid.New().String()[:8]
}
