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
	isCarParkedAnywhere(car *models.Car) bool
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
	if a.isCarParkedAnywhere(car) {
		return nil, errors.ErrCarAlreadyParked
	}

	for _, lot := range a.ParkingLots {
		if !lot.IsFull() {
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

func (a *ParkingAttendant) isCarParkedAnywhere(car *models.Car) bool {
	for _, lot := range a.ParkingLots {
		for _, plateNumber := range lot.ParkedCars {
			if plateNumber == car.LicensePlate {
				return true
			}
		}
	}

	return false
}
