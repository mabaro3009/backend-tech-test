package ride_test

import (
	"testing"
	"time"

	"reby/domain/money"
	"reby/domain/ride"
	"reby/pkg/timenow"

	"github.com/stretchr/testify/assert"
)

func TestBasePriceCalculator_Calculate(t *testing.T) {
	tm := timenow.NewFixedTime(time.Now())
	calculator := ride.NewBasePriceCalculator(100, 10, tm)
	t1 := tm.Now().Add(-5 * time.Minute)
	testCases := []struct {
		description   string
		createdAt     time.Time
		finishedAt    *time.Time
		expectedError error
		expectedPrice money.Money
	}{
		{
			description:   "already finished",
			createdAt:     tm.Now().Add(-10 * time.Minute),
			finishedAt:    &t1,
			expectedError: nil,
			expectedPrice: money.NewMoney(150, "EUR"),
		},
		{
			description:   "invalid timenow",
			createdAt:     tm.Now().Add(10 * time.Minute),
			finishedAt:    &t1,
			expectedError: ride.ErrInvalidRideMinutes,
			expectedPrice: money.NewMoney(0, ""),
		},
		{
			description:   "not finished",
			createdAt:     tm.Now().Add(-10 * time.Minute),
			finishedAt:    nil,
			expectedError: nil,
			expectedPrice: money.NewMoney(200, "EUR"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			r := ride.Ride{
				ID:         "a",
				VehicleID:  "b",
				UserID:     "c",
				StartedAt:  tc.createdAt,
				FinishedAt: tc.finishedAt,
				Price:      nil,
			}

			price, err := calculator.Calculate(r)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, tc.expectedPrice, price)
		})
	}
}
