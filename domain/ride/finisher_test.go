package ride_test

import (
	"context"
	"testing"
	"time"

	"reby/domain/money"
	"reby/domain/ride"
	"reby/pkg/timenow"

	"github.com/stretchr/testify/assert"
)

func TestFinish(t *testing.T) {
	var rideRepoMock *ride.RepoMock
	var priceMock *ride.PriceCalculatorMock
	fixedTime := timenow.NewFixedTime(time.Now())
	now := fixedTime.Now()
	var finisher ride.Finisher
	rideID := "r_1"
	price := money.NewMoney(200, "EUR")
	ctx := context.Background()

	setup := func() {
		rideRepoMock = ride.NewRepoMock()
		priceMock = ride.NewPriceCalculatorMock()
		finisher = ride.NewFinisher(rideRepoMock, priceMock, fixedTime)
	}

	testCases := []struct {
		description string
		finishedAt  *time.Time
		expectedErr error
	}{
		{
			description: "ok",
			finishedAt:  nil,
			expectedErr: nil,
		},
		{
			description: "already finished",
			finishedAt:  &now,
			expectedErr: ride.ErrAlreadyFinished,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			setup()
			startedRide := &ride.Ride{
				ID:         rideID,
				VehicleID:  "1",
				UserID:     "1",
				StartedAt:  now.Add(-5 * time.Minute),
				FinishedAt: tc.finishedAt,
				Price:      nil,
			}

			rideRepoMock.On("GetByID", rideID).Return(startedRide, nil)
			priceMock.On("Calculate", *startedRide).Return(price, nil)

			finishedRide := *startedRide
			finishedRide.FinishedAt = &now
			finishedRide.Price = &price

			rideRepoMock.On("Update", &finishedRide).Return(&finishedRide, nil)

			r, err := finisher.Finish(ctx, rideID)
			assert.ErrorIs(t, tc.expectedErr, err)
			if err == nil {
				assert.NotEmpty(t, r)
			}
		})
	}
}
