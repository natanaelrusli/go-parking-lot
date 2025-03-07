package parkinglot

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/ticket"
)

type ParkingLot struct {
	*models.ParkingLot
}

func New(capacity int) *ParkingLot {
	return &ParkingLot{
		ParkingLot: &models.ParkingLot{
			ID:          uuid.New().String()[:8],
			ParkedCars:  make(map[string]string),
			UsedTickets: make(map[string]bool),
			Capacity:    capacity,
			Subscribers: []models.ParkingLotObserver{},
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

func (p *ParkingLot) IsFull() bool {
	return len(p.ParkedCars) >= p.Capacity
}

// Adding new observers
// why we can use interface as the observer?
func (p *ParkingLot) AddObserver(observer models.ParkingLotObserver) {
	p.Subscribers = append(p.Subscribers, observer)
}

// Notifying all observers
func (p *ParkingLot) notifyObservers() {
	status := models.ParkingLotStatus{
		IsFull:    len(p.ParkedCars) >= p.Capacity,
		LotID:     p.ID,
		Capacity:  p.Capacity,
		Available: p.Capacity - len(p.ParkedCars),
	}

	if status.IsFull {
		fmt.Printf("ALERT: Parking lot %s is now FULL (Capacity: %d)\n",
			status.LotID, status.Capacity)
	} else {
		fmt.Printf("Parking lot %s has %d available spaces (Capacity: %d)\n",
			status.LotID, status.Available, status.Capacity)
	}

	// Notify each observer
	for _, observer := range p.Subscribers {
		observer.OnParkingLotStatusChanged(status)
	}
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

	ticketNumber := ticket.GenerateTicketNumber()
	p.ParkedCars[ticketNumber] = car.LicensePlate

	// Notify observers after successful parking
	p.notifyObservers()

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

	// Notify observers after successful unparking
	p.notifyObservers()

	return car, nil
}
