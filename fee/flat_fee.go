package fee

import "time"

type FlatFeeStrategy struct {
	flatFee float64
}

func NewFlatFeeStrategy(flatFee float64) ParkingFeeStrategy {
	return &FlatFeeStrategy{flatFee: flatFee}
}

func (s *FlatFeeStrategy) CalculateFee(duration time.Duration) float64 {
	return s.flatFee
}
