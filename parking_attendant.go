package main

import (
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
)

type ParkingAttendant struct {
	Name        string
	ParkingLots []*ParkingLot
}

func NewParkingAttendant(name string, parkingLots []*ParkingLot) *ParkingAttendant {
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
		if len(lot.ParkedCars) < lot.Capacity {
			return lot.Park(car)
		}
	}

	return nil, errors.ErrAllLotsAreFull
}

func (a *ParkingAttendant) UnparkCar(ticket *models.Ticket) (*models.Car, error) {
	return a.ParkingLots[0].Unpark(ticket)
}
