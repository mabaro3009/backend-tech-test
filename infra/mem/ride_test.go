package mem_test

import (
	"context"
	"testing"
	"time"

	"reby/domain/money"
	"reby/domain/ride"
	"reby/infra/mem"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRideGetByID(t *testing.T) {
	db := mem.NewRideDB()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		r := &ride.Ride{ID: "1"}
		require.NoError(t, db.Create(ctx, r))

		r, err := db.GetByID(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "1", r.ID)
	})

	t.Run("error", func(t *testing.T) {
		r, err := db.GetByID(ctx, "10")
		assert.ErrorIs(t, ride.ErrNotFound, err)
		assert.Nil(t, r)
	})
}

func TestCreate(t *testing.T) {
	db := mem.NewRideDB()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		r := &ride.Ride{ID: "1"}
		require.NoError(t, db.Create(ctx, r))
	})

	t.Run("already exists", func(t *testing.T) {
		r := &ride.Ride{ID: "2"}
		require.NoError(t, db.Create(ctx, r))

		assert.ErrorIs(t, db.Create(ctx, r), ride.ErrAlreadyExists)
	})
}

func TestIsUserRiding(t *testing.T) {
	db := mem.NewRideDB()
	ctx := context.Background()

	t.Run("no rides from user", func(t *testing.T) {
		userID := "test_user_1"
		isRiding, err := db.IsUserRiding(ctx, userID)
		require.NoError(t, err)
		assert.False(t, isRiding)
	})

	t.Run("user is riding", func(t *testing.T) {
		userID := "test_user_2"
		r := &ride.Ride{
			ID:         "1",
			VehicleID:  "1",
			UserID:     userID,
			StartedAt:  time.Now(),
			FinishedAt: nil,
			Price:      nil,
		}
		require.NoError(t, db.Create(ctx, r))

		isRiding, err := db.IsUserRiding(ctx, userID)
		require.NoError(t, err)
		assert.True(t, isRiding)
	})

	t.Run("user has finished ride", func(t *testing.T) {
		userID := "test_user_3"
		now := time.Now()
		r := &ride.Ride{
			ID:         "2",
			VehicleID:  "2",
			UserID:     userID,
			StartedAt:  now,
			FinishedAt: &now,
			Price:      nil,
		}
		require.NoError(t, db.Create(ctx, r))

		isRiding, err := db.IsUserRiding(ctx, userID)
		require.NoError(t, err)
		assert.False(t, isRiding)
	})
}

func TestIsVehicleRiding(t *testing.T) {
	db := mem.NewRideDB()
	ctx := context.Background()

	t.Run("no rides from vehicle", func(t *testing.T) {
		vID := "test_vehicle_1"
		isRiding, err := db.IsVehicleRiding(ctx, vID)
		require.NoError(t, err)
		assert.False(t, isRiding)
	})

	t.Run("vehicle is riding", func(t *testing.T) {
		vID := "test_vehicle_2"
		r := &ride.Ride{
			ID:         "1",
			VehicleID:  vID,
			UserID:     "1",
			StartedAt:  time.Now(),
			FinishedAt: nil,
			Price:      nil,
		}
		require.NoError(t, db.Create(ctx, r))

		isRiding, err := db.IsVehicleRiding(ctx, vID)
		require.NoError(t, err)
		assert.True(t, isRiding)
	})

	t.Run("vehicle has finished ride", func(t *testing.T) {
		vID := "test_vehicle_3"
		now := time.Now()
		r := &ride.Ride{
			ID:         "2",
			VehicleID:  vID,
			UserID:     "2",
			StartedAt:  now,
			FinishedAt: &now,
			Price:      nil,
		}
		require.NoError(t, db.Create(ctx, r))

		isRiding, err := db.IsVehicleRiding(ctx, vID)
		require.NoError(t, err)
		assert.False(t, isRiding)
	})
}

func TestRideUpdate(t *testing.T) {
	db := mem.NewRideDB()
	ctx := context.Background()

	t.Run("does not exist", func(t *testing.T) {
		r := &ride.Ride{
			ID:         "1",
			VehicleID:  "1",
			UserID:     "1",
			StartedAt:  time.Now(),
			FinishedAt: nil,
			Price:      nil,
		}
		newRide, err := db.Update(ctx, r)
		assert.ErrorIs(t, ride.ErrNotFound, err)
		assert.Nil(t, newRide)
	})

	t.Run("update all fields", func(t *testing.T) {
		now := time.Now()
		r := &ride.Ride{
			ID:         "1",
			VehicleID:  "1",
			UserID:     "1",
			StartedAt:  now,
			FinishedAt: nil,
			Price:      nil,
		}
		require.NoError(t, db.Create(ctx, r))

		updatedRide := &ride.Ride{
			ID:         "1",
			VehicleID:  "2",
			UserID:     "2",
			StartedAt:  now.Add(-5 * time.Minute),
			FinishedAt: &now,
			Price: &money.Money{
				Value:    100,
				Currency: "EUR",
			},
		}

		newRide, err := db.Update(ctx, updatedRide)
		require.NoError(t, err)
		assert.Equal(t, r.ID, newRide.ID)
		assert.Equal(t, r.VehicleID, newRide.VehicleID)
		assert.Equal(t, r.UserID, newRide.UserID)
		assert.Equal(t, r.StartedAt, newRide.StartedAt)
		assert.Equal(t, updatedRide.FinishedAt, newRide.FinishedAt)
		assert.Equal(t, updatedRide.Price, newRide.Price)
	})
}
