package errors

import (
	"errors"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrNilCar message",
			err:      ErrNilCar,
			expected: "cannot park nil car",
		},
		{
			name:     "ErrEmptyLicensePlate message",
			err:      ErrEmptyLicensePlate,
			expected: "cannot park without license plate",
		},
		{
			name:     "ErrNoAvailablePosition message",
			err:      ErrNoAvailablePosition,
			expected: "no available position",
		},
		{
			name:     "ErrCarAlreadyParked message",
			err:      ErrCarAlreadyParked,
			expected: "car already parked",
		},
		{
			name:     "ErrNilTicket message",
			err:      ErrNilTicket,
			expected: "cannot unpark without ticket",
		},
		{
			name:     "ErrEmptyTicketNumber message",
			err:      ErrEmptyTicketNumber,
			expected: "cannot unpark without ticket number",
		},
		{
			name:     "ErrUnrecognizedTicket message",
			err:      ErrUnrecognizedTicket,
			expected: "unrecognized parking ticket",
		},
		{
			name:     "ErrAllLotsAreFull message",
			err:      ErrAllLotsAreFull,
			expected: "all parking lots are full",
		},
		{
			name:     "ErrTicketNotFound message",
			err:      ErrTicketNotFound,
			expected: "ticket not found in any parking lot",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.expected {
				t.Errorf("Expected error message %q, got %q", tt.expected, tt.err.Error())
			}
		})
	}
}

func TestErrorComparisons(t *testing.T) {
	tests := []struct {
		name     string
		err1     error
		err2     error
		expected bool
	}{
		{
			name:     "Compare ErrNilCar with same error",
			err1:     ErrNilCar,
			err2:     ErrNilCar,
			expected: true,
		},
		{
			name:     "Compare ErrNilCar with different error",
			err1:     ErrNilCar,
			err2:     ErrEmptyLicensePlate,
			expected: false,
		},
		{
			name:     "Compare ErrNilCar with new error with same message",
			err1:     ErrNilCar,
			err2:     errors.New("cannot park nil car"),
			expected: true,
		},
		{
			name:     "Compare ErrAllLotsAreFull with same error",
			err1:     ErrAllLotsAreFull,
			err2:     ErrAllLotsAreFull,
			expected: true,
		},
		{
			name:     "Compare ErrTicketNotFound with different error",
			err1:     ErrTicketNotFound,
			err2:     ErrAllLotsAreFull,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err1.Error() == tt.err2.Error()
			if result != tt.expected {
				t.Errorf("Expected comparison to be %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestErrorsAreDistinct(t *testing.T) {
	// Create a map to store all error messages
	errorMessages := make(map[string]error)

	// List of all errors to check
	errors := []error{
		ErrNilCar,
		ErrEmptyLicensePlate,
		ErrNoAvailablePosition,
		ErrCarAlreadyParked,
		ErrNilTicket,
		ErrEmptyTicketNumber,
		ErrUnrecognizedTicket,
		ErrAllLotsAreFull,
		ErrTicketNotFound,
	}

	// Check for duplicate error messages
	for _, err := range errors {
		msg := err.Error()
		if existing, exists := errorMessages[msg]; exists {
			t.Errorf("Duplicate error message found: %q used by multiple errors: %v and %v",
				msg, existing, err)
		}
		errorMessages[msg] = err
	}
}
