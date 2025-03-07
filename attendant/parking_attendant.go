package attendant

import (
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/parkinglot"
)

type ParkingAttendant struct {
	Name          string
	ParkingLots   []*parkinglot.ParkingLot
	AvailableLots map[string]bool
	ParkingStyle  string
}

type ParkingAttendantItf interface {
	GetName() string
	ParkCar(car *models.Car) (*models.Ticket, error)
	UnparkCar(ticket *models.Ticket) (*models.Car, error)
	isCarParkedAnywhere(car *models.Car) bool
	GetAvailableLotsLen() int
	OnParkingLotStatusChanged(status models.ParkingLotStatus)
	GetAllAvailableLots() map[string]bool
	AssignParkingLot(lot *parkinglot.ParkingLot)
	SetParkingStrategy(strategy string)
}

func NewParkingAttendant(name string, parkingLots []*parkinglot.ParkingLot) ParkingAttendantItf {
	availableLots := make(map[string]bool)
	for _, v := range parkingLots {
		availableLots[v.ID] = true
	}

	return &ParkingAttendant{
		Name:          name,
		ParkingLots:   parkingLots,
		AvailableLots: availableLots,
	}
}

func (a *ParkingAttendant) AssignParkingLot(lot *parkinglot.ParkingLot) {
	a.ParkingLots = append(a.ParkingLots, lot)
}

func (a *ParkingAttendant) SetParkingStrategy(strategy string) {
	a.ParkingStyle = strategy
}

func (a *ParkingAttendant) OnParkingLotStatusChanged(status models.ParkingLotStatus) {
	if status.IsFull {
		delete(a.AvailableLots, status.LotID)
		return
	}

	if !status.IsFull && status.Available == 1 {
		a.AvailableLots[status.LotID] = true
		return
	}
}

func (a *ParkingAttendant) GetName() string {
	return a.Name
}

func (a *ParkingAttendant) GetAvailableLotsLen() int {
	return len(a.AvailableLots)
}

func (a *ParkingAttendant) GetAllAvailableLots() map[string]bool {
	return a.AvailableLots
}

func parkInHighestCapacity(parkingLots []*parkinglot.ParkingLot) *parkinglot.ParkingLot {
	hisghestCapacity := 0
	lotWithHighestCapacity := &parkinglot.ParkingLot{}

	for _, v := range parkingLots {
		if v.Capacity > hisghestCapacity {
			hisghestCapacity = v.Capacity
			lotWithHighestCapacity = v
		}
	}

	return lotWithHighestCapacity
}

func parkInHighestFreeSpace(parkingLots []*parkinglot.ParkingLot) *parkinglot.ParkingLot {
	highestFreeSpace := 0
	lotWithHighestFreeSpace := &parkinglot.ParkingLot{}

	for _, v := range parkingLots {
		freeSpace := v.Capacity - v.GetParkedCarCount()

		if v.Capacity > highestFreeSpace {
			highestFreeSpace = freeSpace
			lotWithHighestFreeSpace = v
		}
	}

	return lotWithHighestFreeSpace
}

func (a *ParkingAttendant) ParkCar(car *models.Car) (*models.Ticket, error) {
	if a.isCarParkedAnywhere(car) {
		return nil, errors.ErrCarAlreadyParked
	}

	// with parking style
	if a.ParkingStyle == "capacity" {
		return parkInHighestCapacity(a.ParkingLots).Park(car)
	} else if a.ParkingStyle == "space" {
		return parkInHighestFreeSpace(a.ParkingLots).Park(car)
	}

	// without parking style
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
