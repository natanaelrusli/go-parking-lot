package parking_styles

import (
	"testing"

	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/parkinglot"
	"github.com/stretchr/testify/assert"
)

func TestParkingStyles(t *testing.T) {
	t.Run("should choose lot with highest capacity", func(t *testing.T) {
		cs := NewMostCapacityStrategy()
		pl1 := parkinglot.New(1)
		pl2 := parkinglot.New(2)

		choosenLot := cs.GetLot([]*parkinglot.ParkingLot{
			pl1.(*parkinglot.ParkingLot),
			pl2.(*parkinglot.ParkingLot),
		})

		assert.Equal(t, choosenLot, pl2)
	})

	t.Run("should choose lot with highest free space", func(t *testing.T) {
		cs := NewMostFreeSpaceStrategy()
		pl1 := parkinglot.New(3)
		pl2 := parkinglot.New(2)

		c1 := car.NewCar("B5647PPP")
		c2 := car.NewCar("B5647PPO")

		pl1.Park(c1)
		pl1.Park(c2)

		choosenLot := cs.GetLot([]*parkinglot.ParkingLot{
			pl1.(*parkinglot.ParkingLot),
			pl2.(*parkinglot.ParkingLot),
		})

		assert.Equal(t, choosenLot, pl2)
	})
}
