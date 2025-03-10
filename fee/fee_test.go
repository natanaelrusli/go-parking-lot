package fee_test

import (
	"testing"
	"time"

	"github.com/natanaelrusli/parking-lot/fee"
	"github.com/natanaelrusli/parking-lot/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewFlatFeeStrategy(t *testing.T) {
	t.Run("should return a flat fee strategy", func(t *testing.T) {
		strategy := fee.NewFlatFeeStrategy(10.0)

		assert.NotNil(t, strategy)
	})

	t.Run("mock fee strategy should be called with correct duration", func(t *testing.T) {
		// Arrange
		mockStrategy := mocks.NewParkingFeeStrategy(t)
		duration := 2 * time.Hour

		// Setup expectation
		mockStrategy.On("CalculateFee", duration).Return(20.0).Once()

		// Act
		fee := mockStrategy.CalculateFee(duration)

		// Assert
		assert.Equal(t, 20.0, fee)
		mockStrategy.AssertExpectations(t)
	})
}
