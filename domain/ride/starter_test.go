package ride_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"reby/domain/ride"
	"reby/domain/user"
	"reby/domain/vehicle"
	"reby/pkg/id"
	"reby/pkg/timenow"
)

func TestStart(t *testing.T) {
	var userRepoMock *user.RepoMock
	var vehicleRepoMock *vehicle.RepoMock
	var rideRepoMock *ride.RepoMock
	var idGenMock *id.GeneratorMock
	var starter ride.Starter

	now := time.Now()
	ctx := context.Background()
	setup := func() {
		userRepoMock = user.NewRepoMock()
		vehicleRepoMock = vehicle.NewRepoMock()
		rideRepoMock = ride.NewRepoMock()
		idGenMock = id.NewGeneratorMock()

		fixedTime := timenow.NewFixedTime(now)

		starter = ride.NewStarter(userRepoMock, vehicleRepoMock, rideRepoMock, idGenMock, fixedTime)
	}
	testCases := []struct {
		description     string
		isUserRiding    bool
		isVehicleRiding bool
		expectedError   error
	}{
		{
			description:     "user already riding",
			isUserRiding:    true,
			isVehicleRiding: false,
			expectedError:   ride.ErrUserIsRiding,
		},
		{
			description:     "vehicle already riding",
			isUserRiding:    false,
			isVehicleRiding: true,
			expectedError:   ride.ErrVehicleIsRiding,
		},
		{
			description:     "ok",
			isUserRiding:    false,
			isVehicleRiding: false,
			expectedError:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			setup()

			userID := "u_1"
			vehicleID := "v_1"
			rideID := "r_1"

			userRepoMock.On("GetByID", userID).Return(&user.User{ID: userID}, nil)
			vehicleRepoMock.On("GetByID", vehicleID).Return(&vehicle.Vehicle{ID: vehicleID}, nil)
			rideRepoMock.On("IsUserRiding", userID).Return(tc.isUserRiding, nil)
			rideRepoMock.On("IsVehicleRiding", vehicleID).Return(tc.isVehicleRiding, nil)
			idGenMock.On("Generate").Return(rideID)

			r := &ride.Ride{
				ID:         rideID,
				VehicleID:  vehicleID,
				UserID:     userID,
				StartedAt:  now,
				FinishedAt: nil,
				Price:      nil,
			}

			rideRepoMock.On("Create", r).Return(nil)

			_, err := starter.Start(ctx, ride.StartParams{
				UserID:    userID,
				VehicleID: vehicleID,
			})

			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
