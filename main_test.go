package main

import (
	"testing"
)

func TestGenerateTicketNumber(t *testing.T) {
	// Test that multiple ticket numbers are unique
	tickets := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		ticketNum := generateTicketNumber()

		// Check length
		if len(ticketNum) != 8 {
			t.Errorf("Ticket number length should be 8, got %d", len(ticketNum))
		}

		// Check uniqueness
		if tickets[ticketNum] {
			t.Errorf("Duplicate ticket number generated: %s", ticketNum)
		}
		tickets[ticketNum] = true
	}
}

func TestParkingLotOperations(t *testing.T) {
	parkingLot := NewParkingLot(10)

	// Test parking a car
	t.Run("Park car", func(t *testing.T) {
		car := NewCar("ABC123")
		ticket, _ := parkingLot.Park(car)

		if ticket == nil {
			t.Error("Expected ticket, got nil")
			return
		}
		if _, exists := parkingLot.parkedCars[ticket.ticketNumber]; !exists {
			t.Error("Car should be in parked cars map")
		}
	})

	// Test unparking a car
	t.Run("Unpark car", func(t *testing.T) {
		car := NewCar("XYZ789")
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

		if unparkedCar.licensePlate != "XYZ789" {
			t.Errorf("Expected license plate XYZ789, got %s", unparkedCar.licensePlate)
			return
		}

		if _, exists := parkingLot.parkedCars[ticket.ticketNumber]; exists {
			t.Error("Car should not be in parked cars map after unparking")
		}
	})

	// Test unparking with invalid ticket
	t.Run("Unpark with invalid ticket", func(t *testing.T) {
		invalidTicket := &Ticket{ticketNumber: "INVALID"}
		unparkedCar, err := parkingLot.Unpark(invalidTicket)

		if err == nil {
			t.Error("Expected error for invalid ticket")
		}

		if unparkedCar != nil {
			t.Error("Expected nil car for invalid ticket")
		}
	})

	t.Run("Unpark with used ticket", func(t *testing.T) {
		car := NewCar("ABC123ABC")
		ticket, _ := parkingLot.Park(car)

		_, _ = parkingLot.Unpark(ticket)
		_, err := parkingLot.Unpark(ticket)

		if err == nil {
			t.Error("Expected error for used ticket")
		}

		if err.Error() != "unrecognized parking ticket" {
			t.Error("Expected error for used ticket")
		}
	})

	// Test multiple parking operations
	t.Run("Multiple parking operations", func(t *testing.T) {
		parkingLot := NewParkingLot(10)

		// Park multiple cars
		car1 := NewCar("AAA111")
		car2 := NewCar("BBB222")
		car3 := NewCar("CCC333")

		ticket1, _ := parkingLot.Park(car1)
		ticket2, _ := parkingLot.Park(car2)
		ticket3, _ := parkingLot.Park(car3)

		// Verify all tickets are unique
		if ticket1.ticketNumber == ticket2.ticketNumber ||
			ticket2.ticketNumber == ticket3.ticketNumber ||
			ticket1.ticketNumber == ticket3.ticketNumber {
			t.Error("Ticket numbers should be unique")
		}

		// Verify all cars are parked
		if len(parkingLot.parkedCars) != 3 {
			t.Errorf("Expected 3 parked cars, got %d", len(parkingLot.parkedCars))
		}

		// Unpark cars and verify
		car1Retrieved, _ := parkingLot.Unpark(ticket1)
		if car1Retrieved.licensePlate != "AAA111" {
			t.Errorf("Expected AAA111, got %s", car1Retrieved.licensePlate)
		}

		car2Retrieved, _ := parkingLot.Unpark(ticket2)
		if car2Retrieved.licensePlate != "BBB222" {
			t.Errorf("Expected BBB222, got %s", car2Retrieved.licensePlate)
		}

		// Verify remaining parked cars
		if len(parkingLot.parkedCars) != 1 {
			t.Errorf("Expected 1 parked car, got %d", len(parkingLot.parkedCars))
		}
	})

	t.Run("should return error when parking the same car twice", func(t *testing.T) {
		parkingLot := NewParkingLot(10)

		car := NewCar("AAA111")
		_, _ = parkingLot.Park(car)
		_, err := parkingLot.Park(car)

		if err == nil {
			t.Error("Expected error when parking the same car twice")
		}
	})
}

func TestParkingLotCapacity(t *testing.T) {
	t.Run("should return error when parking 3 cars in a 2 car capacity parking lot", func(t *testing.T) {
		parkingLot := NewParkingLot(2)

		// Park two cars
		car1 := NewCar("AAA111")
		car2 := NewCar("BBB222")
		_, _ = parkingLot.Park(car1)
		_, _ = parkingLot.Park(car2)

		// Try parking a third car, should fail
		car3 := NewCar("CCC333")
		_, err := parkingLot.Park(car3)

		if err == nil {
			t.Error("Expected error when parking over capacity")
		}

		if len(parkingLot.parkedCars) != 2 {
			t.Errorf("Expected 2 parked cars, got %d", len(parkingLot.parkedCars))
		}
	})
}
