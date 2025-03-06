package main

import (
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
)

type ParkingLot struct {
	*models.ParkingLot
}

func NewParkingLot(capacity int) *ParkingLot {
	return &ParkingLot{
		ParkingLot: &models.ParkingLot{
			ParkedCars:  make(map[string]string),
			UsedTickets: make(map[string]bool),
			Capacity:    capacity,
		},
	}
}

func (p *ParkingLot) checkCarExist(car *models.Car) bool {
	for _, plateNumber := range p.ParkedCars {
		if plateNumber == car.LicensePlate {
			return true
		}
	}
	return false
}

func (p *ParkingLot) Park(car *models.Car) (*models.Ticket, error) {
	if car == nil {
		return nil, errors.ErrNilCar
	}

	if car.LicensePlate == "" {
		return nil, errors.ErrEmptyLicensePlate
	}

	if len(p.ParkedCars) >= p.Capacity {
		return nil, errors.ErrNoAvailablePosition
	}

	if p.checkCarExist(car) {
		return nil, errors.ErrCarAlreadyParked
	}

	ticketNumber := generateTicketNumber()
	p.ParkedCars[ticketNumber] = car.LicensePlate

	return &models.Ticket{
		TicketNumber: ticketNumber,
	}, nil
}

func (p *ParkingLot) GetCapacity() int {
	return p.Capacity
}

func (p *ParkingLot) GetParkedCars(ticket *models.Ticket) *models.Car {
	licensePlate, exists := p.ParkedCars[ticket.TicketNumber]
	if !exists {
		return nil
	}

	return &models.Car{
		LicensePlate: licensePlate,
	}
}

func (p *ParkingLot) Unpark(ticket *models.Ticket) (*models.Car, error) {
	if ticket == nil {
		return nil, errors.ErrNilTicket
	}

	if ticket.TicketNumber == "" {
		return nil, errors.ErrEmptyTicketNumber
	}

	if p.UsedTickets[ticket.TicketNumber] {
		return nil, errors.ErrUnrecognizedTicket
	}

	car := p.GetParkedCars(ticket)
	if car == nil {
		return nil, errors.ErrUnrecognizedTicket
	}

	delete(p.ParkedCars, ticket.TicketNumber)
	p.UsedTickets[ticket.TicketNumber] = true

	return car, nil
}
