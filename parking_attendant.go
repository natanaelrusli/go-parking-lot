package main

func NewParkingAttendant(name string, parkingLot *ParkingLot) *ParkingAttendant {
	return &ParkingAttendant{
		name:       name,
		parkingLot: parkingLot,
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
