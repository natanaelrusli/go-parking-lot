package main

import (
	"testing"

	"github.com/natanaelrusli/parking-lot/car"
	"github.com/natanaelrusli/parking-lot/errors"
	"github.com/natanaelrusli/parking-lot/models"
	"github.com/natanaelrusli/parking-lot/parkinglot"
	"github.com/natanaelrusli/parking-lot/ticket"
)

func TestGenerateTicketNumber(t *testing.T) {
	tickets := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		ticketNum := ticket.GenerateTicketNumber()

		if len(ticketNum) != 8 {
			t.Errorf("Ticket number length should be 8, got %d", len(ticketNum))
		}

		if tickets[ticketNum] {
			t.Errorf("Duplicate ticket number generated: %s", ticketNum)
		}
		tickets[ticketNum] = true
	}
}

func TestParkingLotOperations(t *testing.T) {
	parkingLot := parkinglot.New(10)

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

	// Test unparking with invalid ticket
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

	// Test multiple parking operations
	t.Run("Multiple parking operations", func(t *testing.T) {
		parkingLot := parkinglot.New(10)

		// Park multiple cars
		car1 := car.NewCar("AAA111")
		car2 := car.NewCar("BBB222")
		car3 := car.NewCar("CCC333")

		ticket1, _ := parkingLot.Park(car1)
		ticket2, _ := parkingLot.Park(car2)
		ticket3, _ := parkingLot.Park(car3)

		// Verify all tickets are unique
		if ticket1.TicketNumber == ticket2.TicketNumber ||
			ticket2.TicketNumber == ticket3.TicketNumber ||
			ticket1.TicketNumber == ticket3.TicketNumber {
			t.Error("Ticket numbers should be unique")
		}

		// Unpark cars and verify
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
		parkingLot := parkinglot.New(10)

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
		parkingLot := parkinglot.New(2)

		// Park two cars
		car1 := car.NewCar("AAA111")
		car2 := car.NewCar("BBB222")
		_, _ = parkingLot.Park(car1)
		_, _ = parkingLot.Park(car2)

		// Try parking a third car, should fail
		car3 := car.NewCar("CCC333")
		_, err := parkingLot.Park(car3)

		if err == nil {
			t.Error("Expected error when parking over capacity")
		}
	})
}
