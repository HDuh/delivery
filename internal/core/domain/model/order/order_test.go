package order

import (
	"testing"

	"delivery/internal/core/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validLocation(t *testing.T) model.Location {
	t.Helper()
	loc, err := model.NewLocation(5, 5)
	require.NoError(t, err)
	return loc
}

func validVolume(t *testing.T) model.Volume {
	t.Helper()
	vol, err := model.NewVolume(3)
	require.NoError(t, err)
	return vol
}

func TestNewOrder(t *testing.T) {
	t.Run("valid args returns order with Created status", func(t *testing.T) {
		id := uuid.New()
		loc := validLocation(t)
		vol := validVolume(t)

		o, err := NewOrder(id, loc, vol)
		require.NoError(t, err)
		require.NotNil(t, o)

		assert.Equal(t, id, o.ID())
		assert.Equal(t, loc, o.Location())
		assert.Equal(t, vol, o.Volume())
		assert.Equal(t, StatusCreated, o.Status())
	})

	t.Run("nil id returns error", func(t *testing.T) {
		o, err := NewOrder(uuid.Nil, validLocation(t), validVolume(t))
		assert.Nil(t, o)
		assert.Error(t, err)
	})
}

func TestOrder_Assign(t *testing.T) {
	t.Run("Created order assigns successfully", func(t *testing.T) {
		o, err := NewOrder(uuid.New(), validLocation(t), validVolume(t))
		require.NoError(t, err)

		err = o.Assign()
		assert.NoError(t, err)
		assert.Equal(t, statusAssigned, o.Status())
	})

	t.Run("Assigned order cannot be assigned again", func(t *testing.T) {
		o, err := NewOrder(uuid.New(), validLocation(t), validVolume(t))
		require.NoError(t, err)
		require.NoError(t, o.Assign())

		err = o.Assign()
		assert.ErrorIs(t, err, ErrCannotAssign)
	})

	t.Run("Completed order cannot be assigned", func(t *testing.T) {
		o, err := NewOrder(uuid.New(), validLocation(t), validVolume(t))
		require.NoError(t, err)
		require.NoError(t, o.Assign())
		require.NoError(t, o.Complete())

		err = o.Assign()
		assert.ErrorIs(t, err, ErrCannotAssign)
	})
}

func TestOrder_Complete(t *testing.T) {
	t.Run("Assigned order completes successfully", func(t *testing.T) {
		o, err := NewOrder(uuid.New(), validLocation(t), validVolume(t))
		require.NoError(t, err)
		require.NoError(t, o.Assign())

		err = o.Complete()
		assert.NoError(t, err)
		assert.Equal(t, statusCompleted, o.Status())
	})

	t.Run("Created order cannot be completed", func(t *testing.T) {
		o, err := NewOrder(uuid.New(), validLocation(t), validVolume(t))
		require.NoError(t, err)

		err = o.Complete()
		assert.ErrorIs(t, err, ErrCannotComplete)
	})

	t.Run("Completed order cannot be completed again", func(t *testing.T) {
		o, err := NewOrder(uuid.New(), validLocation(t), validVolume(t))
		require.NoError(t, err)
		require.NoError(t, o.Assign())
		require.NoError(t, o.Complete())

		err = o.Complete()
		assert.ErrorIs(t, err, ErrCannotComplete)
	})
}
