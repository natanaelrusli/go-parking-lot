package attendant

import (
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/parkinglot"
	parking_styles "github.com/natanaelrusli/parking-lot/parkingstyles"
)

type ParkingAttendant struct {
	Name          string
	ParkingLots   []*parkinglot.ParkingLot
	AvailableLots map[string]bool
	ParkingStyle  parking_styles.ParkingStyleStrategy
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
	ChangeParkingStrategy(strategy parking_styles.ParkingStyleStrategy)
	GetParkingLots() []*parkinglot.ParkingLot
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

func (a *ParkingAttendant) ChangeParkingStrategy(strategy parking_styles.ParkingStyleStrategy) {
	a.ParkingStyle = strategy
}

func (a *ParkingAttendant) GetParkingLots() []*parkinglot.ParkingLot {
	return a.ParkingLots
}

func (a *ParkingAttendant) AssignParkingLot(lot *parkinglot.ParkingLot) {
	a.ParkingLots = append(a.ParkingLots, lot)
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

func (a *ParkingAttendant) ParkCar(car *models.Car) (*models.Ticket, error) {
	if a.isCarParkedAnywhere(car) {
		return nil, errors.ErrCarAlreadyParked
	}

	// if a parking style is choosen
	if a.ParkingStyle != nil {
		return a.ParkingStyle.GetLot(a.ParkingLots).Park(car)
	}

	// if no parking style choosen, attendant will prioritize any first lot available
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
