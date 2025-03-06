package car

import "github.com/natanaelrusli/parking-lot/models"

func NewCar(licensePlate string) *models.Car {
	return &models.Car{
		LicensePlate: licensePlate,
	}
}
