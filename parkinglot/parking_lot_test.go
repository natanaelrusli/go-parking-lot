package parkinglot

import (
	"testing"
	"time"

	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/fee"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/stretchr/testify/assert"
)

func TestParkingLotOperations(t *testing.T) {
	parkingLot := New(10)

	t.Run("Park car", func(t *testing.T) {
		car := car.NewCar("ABC123")
		ticket, err := parkingLot.Park(car)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if ticket == nil {
			t.Error("Expected ticket, got nil")
			return
		}
	})

	t.Run("Unpark car", func(t *testing.T) {
		car := car.NewCar("XYZ789")
		ticket, _ := parkingLot.Park(car)

		unparkedCar, err := parkingLot.Unpark(ticket)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
			return
		}

		if unparkedCar == nil {
			t.Error("Expected car, got nil")
			return
		}

		if unparkedCar.LicensePlate != "XYZ789" {
			t.Errorf("Expected license plate XYZ789, got %s", unparkedCar.LicensePlate)
			return
		}
	})

	t.Run("Unpark with invalid ticket", func(t *testing.T) {
		invalidTicket := &models.Ticket{TicketNumber: "INVALID"}
		unparkedCar, err := parkingLot.Unpark(invalidTicket)

		if err == nil {
			t.Error("Expected error for invalid ticket")
		}

		if unparkedCar != nil {
			t.Error("Expected nil car for invalid ticket")
		}
	})

	t.Run("Unpark with used ticket", func(t *testing.T) {
		car := car.NewCar("ABC123ABC")
		ticket, _ := parkingLot.Park(car)

		_, _ = parkingLot.Unpark(ticket)
		_, err := parkingLot.Unpark(ticket)

		if err != errors.ErrUnrecognizedTicket {
			t.Errorf("Expected error %v, got %v", errors.ErrUnrecognizedTicket, err)
		}
	})

	t.Run("Multiple parking operations", func(t *testing.T) {
		parkingLot := New(10)

		car1 := car.NewCar("AAA111")
		car2 := car.NewCar("BBB222")
		car3 := car.NewCar("CCC333")

		ticket1, _ := parkingLot.Park(car1)
		ticket2, _ := parkingLot.Park(car2)
		ticket3, _ := parkingLot.Park(car3)

		if ticket1.TicketNumber == ticket2.TicketNumber ||
			ticket2.TicketNumber == ticket3.TicketNumber ||
			ticket1.TicketNumber == ticket3.TicketNumber {
			t.Error("Ticket numbers should be unique")
		}

		car1Retrieved, _ := parkingLot.Unpark(ticket1)
		if car1Retrieved.LicensePlate != "AAA111" {
			t.Errorf("Expected AAA111, got %s", car1Retrieved.LicensePlate)
		}

		car2Retrieved, _ := parkingLot.Unpark(ticket2)
		if car2Retrieved.LicensePlate != "BBB222" {
			t.Errorf("Expected BBB222, got %s", car2Retrieved.LicensePlate)
		}
	})

	t.Run("should return error when parking the same car twice", func(t *testing.T) {
		parkingLot := New(10)

		car := car.NewCar("AAA111")
		_, _ = parkingLot.Park(car)
		_, err := parkingLot.Park(car)

		if err != errors.ErrCarAlreadyParked {
			t.Errorf("Expected error %v, got %v", errors.ErrCarAlreadyParked, err)
		}
	})
}

func TestParkingLotCapacity(t *testing.T) {
	t.Run("should return error when parking 3 cars in a 2 car capacity parking lot", func(t *testing.T) {
		parkingLot := New(2)

		car1 := car.NewCar("AAA111")
		car2 := car.NewCar("BBB222")
		_, _ = parkingLot.Park(car1)
		_, _ = parkingLot.Park(car2)

		car3 := car.NewCar("CCC333")
		_, err := parkingLot.Park(car3)

		if err == nil {
			t.Error("Expected error when parking over capacity")
		}
	})
}

func TestParkingLotFeeStrategy(t *testing.T) {
	t.Run("should have default hourly fee strategy", func(t *testing.T) {
		parkingLot := New(10)

		fee := parkingLot.CalculateFee(time.Hour * 2)
		if fee != 20.0 {
			t.Errorf("Expected fee to be 20.0, got %f", fee)
		}
	})

	t.Run("should be able to change fee strategy to flat fee strategy", func(t *testing.T) {
		pl := New(10)

		f := pl.CalculateFee(time.Hour * 2)
		if f != 20.0 {
			t.Errorf("Expected fee to be 20.0, got %f", f)
		}

		pl.ChangeFeeStrategy(fee.NewFlatFeeStrategy(200))
		f = pl.CalculateFee(1000)

		assert.Equal(t, f, float64(200))
	})
}
