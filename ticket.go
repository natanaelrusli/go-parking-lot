package main

import "github.com/google/uuid"

func generateTicketNumber() string {
	// Using first 8 characters of UUID
	return uuid.New().String()[:8]
}
