package fee

import "time"

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
