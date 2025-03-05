package main

import (
	"errors"
	"log"

	"github.com/google/uuid"
)

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

func generateTicketNumber() string {
	// Using first 8 characters of UUID
	return uuid.New().String()[:8]
}

func NewParkingAttendant(name string, parkingLot *ParkingLot) *ParkingAttendant {
	return &ParkingAttendant{
		name:       name,
		parkingLot: parkingLot,
	}
}

func NewParkingLot(capacity int) *ParkingLot {
	return &ParkingLot{
		parkedCars:  make(map[string]string),
		usedTickets: make(map[string]bool),
		capacity:    capacity,
	}
}

func NewCar(licensePlate string) *Car {
	return &Car{
		licensePlate: licensePlate,
	}
}

func (a *ParkingAttendant) GetName() string {
	return a.name
}

func (a *ParkingAttendant) ParkCar(car *Car) (*Ticket, error) {
	return a.parkingLot.Park(car)
}

func (a *ParkingAttendant) UnparkCar(ticket *Ticket) (*Car, error) {
	return a.parkingLot.Unpark(ticket)
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

func main() {
	log.Println("Parking lot app!")
}
