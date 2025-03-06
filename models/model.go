package models

type ParkingAttendant struct {
	Name        string
	ParkingLots []*ParkingLot
}

type ParkingLot struct {
	// [ticketNumber]plateNumber
	ParkedCars map[string]string
	// [ticketNumber]bool to mark if a ticket has been used
	UsedTickets map[string]bool
	Capacity    int
}

type Car struct {
	LicensePlate string
}

type Ticket struct {
	TicketNumber string
}
