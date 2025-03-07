package models

import "github.com/natanaelrusli/parking-lot/fee"

type ParkingAttendant struct {
	Name        string
	ParkingLots []*ParkingLot
}

// The status data that gets passed to observers
type ParkingLotStatus struct {
	IsFull    bool
	LotID     string
	Capacity  int
	Available int
}

// The Observer interface
type ParkingLotObserver interface {
	OnParkingLotStatusChanged(status ParkingLotStatus)
}

// The Subject (ParkingLot) containing list of observers
type ParkingLot struct {
	ID          string
	ParkedCars  map[string]string
	UsedTickets map[string]bool
	Capacity    int
	// List of observers
	Subscribers []ParkingLotObserver
	FeeStrategy fee.ParkingFeeStrategy
}

type Car struct {
	LicensePlate string
}

type Ticket struct {
	TicketNumber string
}
