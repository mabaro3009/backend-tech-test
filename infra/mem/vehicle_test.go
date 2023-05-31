package mem_test

import (
	"context"
	"testing"

	"reby/domain/vehicle"
	"reby/infra/mem"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVehicleGetByID(t *testing.T) {
	db := mem.NewVehicleDB()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		v, err := db.GetByID(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "1", v.ID)
	})

	t.Run("not found", func(t *testing.T) {
		v, err := db.GetByID(ctx, "10")
		assert.ErrorIs(t, vehicle.ErrNotFound, err)
		assert.Nil(t, v)
	})
}

func TestCreatVehicle(t *testing.T) {
	db := mem.NewVehicleDB()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		v := &vehicle.Vehicle{ID: "v_test"}
		require.NoError(t, db.Create(ctx, v))
	})

	t.Run("already exists", func(t *testing.T) {
		u := &vehicle.Vehicle{ID: "v_test_2"}
		require.NoError(t, db.Create(ctx, u))
		assert.ErrorIs(t, vehicle.ErrAlreadyExists, db.Create(ctx, u))
	})
}
