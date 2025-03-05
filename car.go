package main

func NewCar(licensePlate string) *Car {
	return &Car{
		licensePlate: licensePlate,
	}
}
