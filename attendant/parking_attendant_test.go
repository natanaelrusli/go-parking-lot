package attendant

import (
	"testing"

	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/parkinglot"
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

func TestParkingStyle(t *testing.T) {
	t.Run("should park in first available lot", func(t *testing.T) {
		pl1 := parkinglot.New(2)
		pl2 := parkinglot.New(3)

		c1 := car.NewCar("B7676POO")

		att := NewParkingAttendant("sule", []*parkinglot.ParkingLot{
			pl1.(*parkinglot.ParkingLot), pl2.(*parkinglot.ParkingLot),
		})

		att.ParkCar(c1)

		count1 := pl1.GetParkedCarCount()
		count2 := pl2.GetParkedCarCount()

		assert.Equal(t, 1, count1)
		assert.Equal(t, 0, count2)
	})

	t.Run("should park in lot with highest capacity", func(t *testing.T) {
		pl1 := parkinglot.New(2)
		pl2 := parkinglot.New(3)

		c1 := car.NewCar("B7676POO")

		att := NewParkingAttendant("sule", []*parkinglot.ParkingLot{
			pl1.(*parkinglot.ParkingLot), pl2.(*parkinglot.ParkingLot),
		})

		att.SetParkingStrategy("capacity")

		att.ParkCar(c1)

		count1 := pl1.GetParkedCarCount()
		count2 := pl2.GetParkedCarCount()

		assert.Equal(t, 0, count1)
		assert.Equal(t, 1, count2)
	})

	t.Run("should park in lot with highest free space", func(t *testing.T) {
		pl1 := parkinglot.New(2)
		pl2 := parkinglot.New(2)

		c0 := car.NewCar("RI1")
		c1 := car.NewCar("B7676POO")

		pl1.Park(c0)

		att := NewParkingAttendant("sule", []*parkinglot.ParkingLot{
			pl1.(*parkinglot.ParkingLot), pl2.(*parkinglot.ParkingLot),
		})

		att.SetParkingStrategy("space")

		att.ParkCar(c1)

		count1 := pl1.GetParkedCarCount()
		count2 := pl2.GetParkedCarCount()

		assert.Equal(t, 1, count1)
		assert.Equal(t, 1, count2)
	})
}
