package main

import (
	"testing"

	"github.com/natanaelrusli/parking-lot/car"
)

func TestParkingAttendant(t *testing.T) {
	t.Run("should be able to park cars", func(t *testing.T) {
		parkingLot := NewParkingLot(10)
		parkingAttendant := NewParkingAttendant("John", []*ParkingLot{parkingLot})

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
		parkingLot := NewParkingLot(10)
		parkingAttendant := NewParkingAttendant("John", []*ParkingLot{parkingLot})

		car := car.NewCar("AAA111")
		ticket, _ := parkingAttendant.ParkCar(car)

		_, err := parkingAttendant.UnparkCar(ticket)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
	})

	t.Run("should park cars in next lot when first lot is full", func(t *testing.T) {
		lot1 := NewParkingLot(1)
		lot2 := NewParkingLot(1)
		attendant := NewParkingAttendant("John", []*ParkingLot{lot1, lot2})

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
		if err == nil {
			t.Error("Expected error when all lots are full")
		}
	})
}
