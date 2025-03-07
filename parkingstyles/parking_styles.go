package parking_styles

import "github.com/natanaelrusli/parking-lot/parkinglot"

type ParkingStyleStrategy interface {
	GetLot(parkingLots []*parkinglot.ParkingLot) *parkinglot.ParkingLot
}

type MostCapacityStrategy struct{}

type MostFreeSpaceStrategy struct{}

func NewMostCapacityStrategy() ParkingStyleStrategy {
	return &MostCapacityStrategy{}
}

func NewMostFreeSpaceStrategy() ParkingStyleStrategy {
	return &MostFreeSpaceStrategy{}
}

func (s *MostCapacityStrategy) GetLot(parkingLots []*parkinglot.ParkingLot) *parkinglot.ParkingLot {
	highestCap := 0
	lotWithHighestCap := &parkinglot.ParkingLot{}

	for _, v := range parkingLots {
		if v.Capacity > highestCap {
			highestCap = v.Capacity
			lotWithHighestCap = v
		}
	}
	return lotWithHighestCap
}

func (s *MostFreeSpaceStrategy) GetLot(parkingLots []*parkinglot.ParkingLot) *parkinglot.ParkingLot {
	highestFree := 0
	lotWithHighestFree := &parkinglot.ParkingLot{}

	for _, v := range parkingLots {
		freeSpace := v.Capacity - len(v.ParkedCars)

		if freeSpace > highestFree {
			highestFree = freeSpace
			lotWithHighestFree = v
		}
	}
	return lotWithHighestFree
}
