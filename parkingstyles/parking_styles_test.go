package parking_styles

import (
	"testing"

	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/parkinglot"
	"github.com/stretchr/testify/assert"
)

func TestMostCapacityStrategy(t *testing.T) {
	t.Run("should choose lot with highest capacity when all lots are empty", func(t *testing.T) {
		// Arrange
		strategy := NewMostCapacityStrategy()
		lot1 := parkinglot.New(3)
		lot2 := parkinglot.New(5)
		lot3 := parkinglot.New(2)

		// Act
		chosenLot := strategy.GetLot([]*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
			lot3.(*parkinglot.ParkingLot),
		})

		// Assert
		assert.Equal(t, lot2, chosenLot)
	})

	t.Run("should choose lot with highest capacity even when some lots are partially filled", func(t *testing.T) {
		// Arrange
		strategy := NewMostCapacityStrategy()
		lot1 := parkinglot.New(3)
		lot2 := parkinglot.New(5)
		lot3 := parkinglot.New(2)

		// Park some cars in lot2
		car1 := car.NewCar("ABC123")
		car2 := car.NewCar("DEF456")
		lot2.Park(car1)
		lot2.Park(car2)

		// Act
		chosenLot := strategy.GetLot([]*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
			lot3.(*parkinglot.ParkingLot),
		})

		// Assert
		assert.Equal(t, lot2, chosenLot)
	})

	t.Run("should choose first lot when multiple lots have same capacity", func(t *testing.T) {
		// Arrange
		strategy := NewMostCapacityStrategy()
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(5)

		// Act
		chosenLot := strategy.GetLot([]*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		// Assert
		assert.Equal(t, lot1, chosenLot)
	})
}

func TestMostFreeSpaceStrategy(t *testing.T) {
	t.Run("should choose lot with most free spaces when lots have different occupancy", func(t *testing.T) {
		// Arrange
		strategy := NewMostFreeSpaceStrategy()
		lot1 := parkinglot.New(3)
		lot2 := parkinglot.New(5)
		lot3 := parkinglot.New(4)

		// Park cars in lots
		car1 := car.NewCar("ABC123")
		car2 := car.NewCar("DEF456")
		car3 := car.NewCar("GHI789")

		lot1.Park(car1) // 2 spaces left
		lot2.Park(car2) // 4 spaces left
		lot2.Park(car3) // 3 spaces left

		// Act
		chosenLot := strategy.GetLot([]*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
			lot3.(*parkinglot.ParkingLot),
		})

		// Assert
		assert.Equal(t, lot3, chosenLot) // lot3 has 4 free spaces
	})

	t.Run("should choose lot with most free spaces when lots have same capacity", func(t *testing.T) {
		// Arrange
		strategy := NewMostFreeSpaceStrategy()
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(5)

		// Park cars in lot1
		car1 := car.NewCar("ABC123")
		car2 := car.NewCar("DEF456")
		lot1.Park(car1)
		lot1.Park(car2)

		// Act
		chosenLot := strategy.GetLot([]*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		// Assert
		assert.Equal(t, lot2, chosenLot) // lot2 has 5 free spaces vs lot1's 3
	})

	t.Run("should choose first lot when multiple lots have same free space", func(t *testing.T) {
		// Arrange
		strategy := NewMostFreeSpaceStrategy()
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(5)

		// Act
		chosenLot := strategy.GetLot([]*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		// Assert
		assert.Equal(t, lot1, chosenLot)
	})

	t.Run("should handle empty lots list", func(t *testing.T) {
		// Arrange
		strategy := NewMostFreeSpaceStrategy()

		// Act
		chosenLot := strategy.GetLot([]*parkinglot.ParkingLot{})

		// Assert
		assert.Equal(t, &parkinglot.ParkingLot{}, chosenLot)
	})
}
