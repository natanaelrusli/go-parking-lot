package main

import "errors"

func NewParkingLot(capacity int) *ParkingLot {
	return &ParkingLot{
		parkedCars:  make(map[string]string),
		usedTickets: make(map[string]bool),
		capacity:    capacity,
	}
}

func (p *ParkingLot) checkCarExist(car *Car) bool {
	for _, plateNumber := range p.parkedCars {
		if plateNumber == car.licensePlate {
			return true
		}
	}

	return false
}

func (p *ParkingLot) Park(car *Car) (*Ticket, error) {
	if car == nil {
		return nil, errors.New("cannot park nil car")
	}

	if car.licensePlate == "" {
		return nil, errors.New("cannot park without license plate")
	}

	if len(p.parkedCars) >= p.capacity {
		return nil, errors.New("no available position")
	}

	// Check if car is already parked by searching through values
	if p.checkCarExist(car) {
		return nil, errors.New("car already parked")
	}

	ticketNumber := generateTicketNumber()
	p.parkedCars[ticketNumber] = car.licensePlate

	return &Ticket{
		ticketNumber: ticketNumber,
	}, nil
}

func (p *ParkingLot) GetCapacity() int {
	return p.capacity
}

func (p *ParkingLot) CheckCarExists(licensePlate string) bool {
	_, exists := p.parkedCars[licensePlate]
	return exists
}

func (p *ParkingLot) GetParkedCars(ticket *Ticket) *Car {
	licensePlate, exists := p.parkedCars[ticket.ticketNumber]
	if !exists {
		return nil
	}

	return &Car{
		licensePlate: licensePlate,
	}
}

func (p *ParkingLot) Unpark(ticket *Ticket) (*Car, error) {
	if ticket == nil {
		return nil, errors.New("cannot unpark without ticket")
	}

	if ticket.ticketNumber == "" {
		return nil, errors.New("cannot unpark without ticket number")
	}

	if p.usedTickets[ticket.ticketNumber] {
		return nil, errors.New("unrecognized parking ticket")
	}

	car := p.GetParkedCars(ticket)

	if car == nil {
		return nil, errors.New("unrecognized parking ticket")
	}

	delete(p.parkedCars, ticket.ticketNumber)
	p.usedTickets[ticket.ticketNumber] = true

	return car, nil
}
