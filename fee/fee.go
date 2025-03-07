package fee

import "time"

type ParkingFeeStrategy interface {
	CalculateFee(duration time.Duration) float64
}

type HourlyFeeStrategy struct {
	ratePerHour float64
}

func NewHourlyFeeStrategy(ratePerHour float64) ParkingFeeStrategy {
	return &HourlyFeeStrategy{ratePerHour: ratePerHour}
}

func (s *HourlyFeeStrategy) CalculateFee(duration time.Duration) float64 {
	hours := float64(duration.Hours())

	if hours < 1 {
		return s.ratePerHour
	}

	return float64(duration.Hours()) * s.ratePerHour
}

type FlatFeeStrategy struct {
	flatFee float64
}

func NewFlatFeeStrategy(flatFee float64) ParkingFeeStrategy {
	return &FlatFeeStrategy{flatFee: flatFee}
}

func (s *FlatFeeStrategy) CalculateFee(duration time.Duration) float64 {
	return s.flatFee
}
