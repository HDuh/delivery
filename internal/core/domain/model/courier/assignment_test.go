package model

import (
	"testing"

	domainmodel "delivery/internal/core/domain/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func validLocation(t *testing.T) domainmodel.Location {
	t.Helper()
	loc, err := domainmodel.NewLocation(5, 5)
	require.NoError(t, err)
	return loc
}

func validVolume(t *testing.T) domainmodel.Volume {
	t.Helper()
	vol, err := domainmodel.NewVolume(3)
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
		a, err := NewAssignment(uuid.New(), domainmodel.QuantityZero, validLocation(t))
		assert.Nil(t, a)
		assert.ErrorIs(t, err, ErrVolumeIsZero)
	})
}

func TestAssignment_Complete(t *testing.T) {
	t.Run("completes assigned assignment without error", func(t *testing.T) {
		a, err := NewAssignment(uuid.New(), validVolume(t), validLocation(t))
		require.NoError(t, err)

		err = a.Complete()
		assert.NoError(t, err)
		assert.True(t, a.Status().Equals(statusCompleted))
	})

	t.Run("completing already completed returns ErrAssignmentAlreadyCompleted", func(t *testing.T) {
		a, err := NewAssignment(uuid.New(), validVolume(t), validLocation(t))
		require.NoError(t, err)
		require.NoError(t, a.Complete())

		err = a.Complete()
		assert.ErrorIs(t, err, ErrAssignmentAlreadyCompleted)
	})
}
