package main

type ParkingAttendant struct {
	name       string
	parkingLot *ParkingLot
}

type ParkingLot struct {
	// [ticketNumber]plateNumber
	parkedCars map[string]string
	// [ticketNumber]bool to mark if a ticket has been used
	usedTickets map[string]bool
	capacity    int
}

type Car struct {
	licensePlate string
}

type Ticket struct {
	ticketNumber string
}
