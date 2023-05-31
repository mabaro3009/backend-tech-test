package mem_test

import (
	"context"
	"testing"

	"reby/domain/user"
	"reby/infra/mem"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserGetByID(t *testing.T) {
	db := mem.NewUserDB()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		u, err := db.GetByID(ctx, "1")
		require.NoError(t, err)
		assert.Equal(t, "1", u.ID)
	})

	t.Run("not found", func(t *testing.T) {
		u, err := db.GetByID(ctx, "10")
		assert.ErrorIs(t, user.ErrNotFound, err)
		assert.Nil(t, u)
	})
}

func TestCreateUser(t *testing.T) {
	db := mem.NewUserDB()
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		u := &user.User{ID: "u_test"}
		require.NoError(t, db.Create(ctx, u))
	})

	t.Run("already exists", func(t *testing.T) {
		u := &user.User{ID: "u_test_2"}
		require.NoError(t, db.Create(ctx, u))
		assert.ErrorIs(t, user.ErrAlreadyExists, db.Create(ctx, u))
	})
}
