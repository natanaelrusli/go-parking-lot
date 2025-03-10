package fee

import "time"

type ParkingFeeStrategy interface {
	CalculateFee(duration time.Duration) float64
}
