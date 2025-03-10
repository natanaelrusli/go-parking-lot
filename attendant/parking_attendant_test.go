package attendant

import (
	"testing"

	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/parkinglot"
	parking_styles "github.com/natanaelrusli/parking-lot/parkingstyles"
	"github.com/stretchr/testify/assert"
)

func TestParkingAttendant(t *testing.T) {
	t.Run("should be able to park cars", func(t *testing.T) {
		parkingLot := parkinglot.New(10)
		parkingAttendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{})
		parkingAttendant.AssignParkingLot(parkingLot.(*parkinglot.ParkingLot))

		car := car.NewCar("AAA111")
		ticket, err := parkingAttendant.ParkCar(car)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if ticket == nil {
			t.Error("Expected ticket, got nil")
		}
	})

	t.Run("should be able to unpark cars", func(t *testing.T) {
		parkingLot := parkinglot.New(10)
		parkingAttendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{parkingLot.(*parkinglot.ParkingLot)})

		car := car.NewCar("AAA111")
		ticket, _ := parkingAttendant.ParkCar(car)

		_, err := parkingAttendant.UnparkCar(ticket)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("should park cars in next lot when first lot is full", func(t *testing.T) {
		lot1 := parkinglot.New(1)
		lot2 := parkinglot.New(1)
		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		car1 := car.NewCar("AAA111")
		car2 := car.NewCar("BBB222")

		// Park first car (should go to first lot)
		ticket1, err := attendant.ParkCar(car1)
		if err != nil {
			t.Errorf("Unexpected error parking first car: %v", err)
		}
		if lot1.GetParkedCars(ticket1) == nil {
			t.Error("First car should be in first lot")
		}

		// Park second car (should go to second lot)
		ticket2, err := attendant.ParkCar(car2)
		if err != nil {
			t.Errorf("Unexpected error parking second car: %v", err)
		}
		if lot2.GetParkedCars(ticket2) == nil {
			t.Error("Second car should be in second lot")
		}

		// Try parking third car (should fail as both lots are full)
		car3 := car.NewCar("CCC333")
		_, err = attendant.ParkCar(car3)
		if err != errors.ErrAllLotsAreFull {
			t.Errorf("Expected error %v when all lots are full, got %v", errors.ErrAllLotsAreFull, err)
		}
	})
}

func TestPreventDoubleParking(t *testing.T) {
	t.Run("should not allow parking the same car twice", func(t *testing.T) {
		lot1 := parkinglot.New(10)
		lot2 := parkinglot.New(10)

		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		car := car.NewCar("AAA111")
		_, _ = attendant.ParkCar(car)
		_, err := attendant.ParkCar(car)

		if err != errors.ErrCarAlreadyParked {
			t.Errorf("Expected error %v, got %v", errors.ErrCarAlreadyParked, err)
		}
	})

	t.Run("should not allow unparking with invalid ticket", func(t *testing.T) {
		lot1 := parkinglot.New(10)
		lot2 := parkinglot.New(10)

		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		invalidTicket := &models.Ticket{TicketNumber: "INVALID"}
		_, err := attendant.UnparkCar(invalidTicket)

		if err != errors.ErrTicketNotFound {
			t.Errorf("Expected error %v, got %v", errors.ErrTicketNotFound, err)
		}
	})

	t.Run("should be able to park in lot with highest cap", func(t *testing.T) {
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(10)
		c1 := car.NewCar("B6565POP")

		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		assert.NotNil(t, attendant)

		ps := parking_styles.NewMostCapacityStrategy()

		attendant.ChangeParkingStrategy(ps)
		ticket, err := attendant.ParkCar(c1)

		assert.NotNil(t, ticket)
		assert.NoError(t, err)

		assert.Equal(t, 0, lot1.GetParkedCarCount())
		assert.Equal(t, 1, lot2.GetParkedCarCount())
	})

	t.Run("should be able to park in lot with highest free space", func(t *testing.T) {
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(5)
		c0 := car.NewCar("RI1")
		c1 := car.NewCar("B6565POP")
		lot1.Park(c0)

		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		assert.NotNil(t, attendant)

		ps := parking_styles.NewMostFreeSpaceStrategy()

		attendant.ChangeParkingStrategy(ps)
		ticket, err := attendant.ParkCar(c1)

		assert.NotNil(t, ticket)
		assert.NoError(t, err)

		assert.Equal(t, 1, lot1.GetParkedCarCount())
		assert.Equal(t, 1, lot2.GetParkedCarCount())
	})
}

func TestAttendantAvailableLots(t *testing.T) {
	t.Run("should have all lots as available initially", func(t *testing.T) {
		// arrange
		pl1 := parkinglot.New(2)
		pl2 := parkinglot.New(2)

		// act
		attendant := NewParkingAttendant("sule", []*parkinglot.ParkingLot{
			pl1.(*parkinglot.ParkingLot),
			pl2.(*parkinglot.ParkingLot),
		})
		al := attendant.GetAvailableLotsLen()

		// assert
		assert.Equal(t, al, 2)
	})

	t.Run("should have remove full lot from availableLots list", func(t *testing.T) {
		// arrange
		pl1 := parkinglot.New(1)
		pl2 := parkinglot.New(2)
		c1 := car.NewCar("B92728POS")
		attendant := NewParkingAttendant("sule", []*parkinglot.ParkingLot{
			pl1.(*parkinglot.ParkingLot),
			pl2.(*parkinglot.ParkingLot),
		})
		pl1.AddObserver(attendant)

		// act
		attendant.ParkCar(c1)
		al := attendant.GetAvailableLotsLen()

		// assert
		assert.Equal(t, al, 1)
	})

	t.Run("should not notify a non subscriber", func(t *testing.T) {
		// arrange
		pl1 := parkinglot.New(1)
		c1 := car.NewCar("B8888POP")
		at1 := NewParkingAttendant("sule", []*parkinglot.ParkingLot{pl1.(*parkinglot.ParkingLot)})
		at2 := NewParkingAttendant("sule", []*parkinglot.ParkingLot{pl1.(*parkinglot.ParkingLot)})
		pl1.AddObserver(at1)

		// act
		at1.ParkCar(c1)
		al1 := at1.GetAvailableLotsLen()
		al2 := at2.GetAvailableLotsLen()

		// assert
		assert.Equal(t, al1, 0)
		assert.Equal(t, al2, 1)
	})

	t.Run("should be notified when a parking lot become available or at least have 1 available space", func(t *testing.T) {
		// arrange
		pl1 := parkinglot.New(1)
		c1 := car.NewCar("B8888POP")
		at1 := NewParkingAttendant("sule", []*parkinglot.ParkingLot{pl1.(*parkinglot.ParkingLot)})
		pl1.AddObserver(at1)

		// act
		ticket1, err := at1.ParkCar(c1)
		al1 := at1.GetAvailableLotsLen()

		// assert
		assert.Equal(t, al1, 0)
		assert.NotNil(t, ticket1)
		assert.NoError(t, err)
		allAvailableLot := at1.GetAllAvailableLots()
		assert.Equal(t, allAvailableLot[pl1.(*parkinglot.ParkingLot).ID], false)

		// act
		returnedCar1, err := at1.UnparkCar(ticket1)
		al1 = at1.GetAvailableLotsLen()

		// assert
		assert.Equal(t, al1, 1)
		assert.NotNil(t, returnedCar1)
		assert.NoError(t, err)

		allAvailableLot = at1.GetAllAvailableLots()
		assert.Equal(t, allAvailableLot[pl1.(*parkinglot.ParkingLot).ID], true)
	})
}

func TestParkingStyles(t *testing.T) {
	t.Run("should park car using default strategy when no style is set", func(t *testing.T) {
		// Arrange
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(10)
		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})
		car := car.NewCar("ABC123")

		// Act
		ticket, err := attendant.ParkCar(car)

		// Assert
		assert.NotNil(t, ticket)
		assert.NoError(t, err)
		assert.Equal(t, 1, lot1.GetParkedCarCount()) // Should use first available lot
		assert.Equal(t, 0, lot2.GetParkedCarCount())
	})

	t.Run("should park car using most capacity strategy", func(t *testing.T) {
		// Arrange
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(10)
		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})
		car := car.NewCar("ABC123")

		// Change to most capacity strategy
		strategy := parking_styles.NewMostCapacityStrategy()
		attendant.ChangeParkingStrategy(strategy)

		// Act
		ticket, err := attendant.ParkCar(car)

		// Assert
		assert.NotNil(t, ticket)
		assert.NoError(t, err)
		assert.Equal(t, 0, lot1.GetParkedCarCount())
		assert.Equal(t, 1, lot2.GetParkedCarCount()) // Should use lot with highest capacity
	})

	t.Run("should park car using most free space strategy", func(t *testing.T) {
		// Arrange
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(5)

		// Park a car in lot1 to make lot2 have more free space
		lot1.Park(car.NewCar("XYZ789"))

		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})
		newCar := car.NewCar("ABC123")

		// Change to most free space strategy
		strategy := parking_styles.NewMostFreeSpaceStrategy()
		attendant.ChangeParkingStrategy(strategy)

		// Act
		ticket, err := attendant.ParkCar(newCar)

		// Assert
		assert.NotNil(t, ticket)
		assert.NoError(t, err)
		assert.Equal(t, 1, lot1.GetParkedCarCount())
		assert.Equal(t, 1, lot2.GetParkedCarCount()) // Should use lot with most free space
	})

	t.Run("should be able to change parking strategy dynamically", func(t *testing.T) {
		// Arrange
		lot1 := parkinglot.New(5)
		lot2 := parkinglot.New(10)
		attendant := NewParkingAttendant("John", []*parkinglot.ParkingLot{
			lot1.(*parkinglot.ParkingLot),
			lot2.(*parkinglot.ParkingLot),
		})

		car1 := car.NewCar("ABC123")
		car2 := car.NewCar("XYZ789")

		// First use most capacity strategy
		attendant.ChangeParkingStrategy(parking_styles.NewMostFreeSpaceStrategy())
		ticket1, _ := attendant.ParkCar(car1)
		assert.Equal(t, 0, lot1.GetParkedCarCount())
		assert.Equal(t, 1, lot2.GetParkedCarCount())

		// Then change to most free space strategy
		attendant.ChangeParkingStrategy(parking_styles.NewMostCapacityStrategy())
		ticket2, _ := attendant.ParkCar(car2)
		assert.Equal(t, 0, lot1.GetParkedCarCount())
		assert.Equal(t, 2, lot2.GetParkedCarCount())

		// Cleanup
		attendant.UnparkCar(ticket1)
		attendant.UnparkCar(ticket2)
	})
}
