package errors

import "errors"

var (
	// Parking lot errors
	ErrNilCar              = errors.New("cannot park nil car")
	ErrEmptyLicensePlate   = errors.New("cannot park without license plate")
	ErrNoAvailablePosition = errors.New("no available position")
	ErrCarAlreadyParked    = errors.New("car already parked")
	ErrNilTicket           = errors.New("cannot unpark without ticket")
	ErrEmptyTicketNumber   = errors.New("cannot unpark without ticket number")
	ErrUnrecognizedTicket  = errors.New("unrecognized parking ticket")

	// Parking attendant errors
	ErrAllLotsAreFull = errors.New("all parking lots are full")
	ErrTicketNotFound = errors.New("ticket not found in any parking lot")
)
