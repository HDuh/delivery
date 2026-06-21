package courier

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

func TestNewAssignment(t *testing.T) {
	t.Run("valid args returns assignment with assigned status", func(t *testing.T) {
		orderID := uuid.New()
		vol := validVolume(t)
		loc := validLocation(t)

		a, err := NewAssignment(orderID, vol, loc)
		require.NoError(t, err)
		require.NotNil(t, a)

		assert.NotEqual(t, uuid.Nil, a.ID())
		assert.Equal(t, orderID, a.OrderID())
		assert.Equal(t, vol, a.Volume())
		assert.Equal(t, loc, a.Location())
		assert.True(t, a.Status().Equals(statusAssigned))
	})

	t.Run("nil orderID returns ErrInvalidOrderID", func(t *testing.T) {
		a, err := NewAssignment(uuid.Nil, validVolume(t), validLocation(t))
		assert.Nil(t, a)
		assert.ErrorIs(t, err, ErrInvalidOrderID)
	})

	t.Run("zero volume returns ErrVolumeIsZero", func(t *testing.T) {
		a, err := NewAssignment(uuid.New(), model.QuantityZero, validLocation(t))
		assert.Nil(t, a)
		assert.ErrorIs(t, err, ErrVolumeIsZero)
	})
}

func TestAssignment_Complete(t *testing.T) {
	t.Run("completes when courier is at same location", func(t *testing.T) {
		loc := validLocation(t)
		a, err := NewAssignment(uuid.New(), validVolume(t), loc)
		require.NoError(t, err)

		err = a.Complete(loc)
		assert.NoError(t, err)
		assert.True(t, a.Status().Equals(statusCompleted))
	})

	t.Run("completes when courier is 1 cell away", func(t *testing.T) {
		orderLoc, err := model.NewLocation(5, 5)
		require.NoError(t, err)
		courierLoc, err := model.NewLocation(5, 6)
		require.NoError(t, err)

		a, err := NewAssignment(uuid.New(), validVolume(t), orderLoc)
		require.NoError(t, err)

		assert.NoError(t, a.Complete(courierLoc))
	})

	t.Run("returns ErrCourierTooFar when courier is more than 1 cell away", func(t *testing.T) {
		orderLoc, err := model.NewLocation(5, 5)
		require.NoError(t, err)
		courierLoc, err := model.NewLocation(5, 7)
		require.NoError(t, err)

		a, err := NewAssignment(uuid.New(), validVolume(t), orderLoc)
		require.NoError(t, err)

		assert.ErrorIs(t, a.Complete(courierLoc), ErrCourierTooFar)
	})

	t.Run("completing already completed returns ErrAssignmentAlreadyCompleted", func(t *testing.T) {
		loc := validLocation(t)
		a, err := NewAssignment(uuid.New(), validVolume(t), loc)
		require.NoError(t, err)
		require.NoError(t, a.Complete(loc))

		err = a.Complete(loc)
		assert.ErrorIs(t, err, ErrAssignmentAlreadyCompleted)
	})
}
