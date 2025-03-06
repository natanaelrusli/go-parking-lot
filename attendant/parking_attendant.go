package attendant

import (
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/parkinglot"
)

type ParkingAttendant struct {
	Name        string
	ParkingLots []*parkinglot.ParkingLot
}

type ParkingAttendantItf interface {
	GetName() string
	ParkCar(car *models.Car) (*models.Ticket, error)
	UnparkCar(ticket *models.Ticket) (*models.Car, error)
}

func NewParkingAttendant(name string, parkingLots []*parkinglot.ParkingLot) ParkingAttendantItf {
	return &ParkingAttendant{
		Name:        name,
		ParkingLots: parkingLots,
	}
}

func (a *ParkingAttendant) GetName() string {
	return a.Name
}

func (a *ParkingAttendant) ParkCar(car *models.Car) (*models.Ticket, error) {
	for _, lot := range a.ParkingLots {
		if len(lot.ParkedCars) < lot.GetCapacity() {
			return lot.Park(car)
		}
	}
	return nil, errors.ErrAllLotsAreFull
}

func (a *ParkingAttendant) UnparkCar(ticket *models.Ticket) (*models.Car, error) {
	for _, lot := range a.ParkingLots {
		if car := lot.GetParkedCars(ticket); car != nil {
			return lot.Unpark(ticket)
		}
	}
	return nil, errors.ErrTicketNotFound
}
