package courier

import (
	"testing"

	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validCourierName = "PiterParker"

func validSpeed(t *testing.T) Speed {
	t.Helper()
	s, err := NewSpeed(2)
	require.NoError(t, err)
	return s
}

func validCourier(t *testing.T) *Courier {
	t.Helper()
	c, err := NewCourier(validCourierName, validSpeed(t), validLocation(t))
	require.NoError(t, err)
	return c
}

func TestNewCourier(t *testing.T) {
	loc := validLocation(t)
	speed := validSpeed(t)

	t.Run("valid args returns courier", func(t *testing.T) {
		c, err := NewCourier(validCourierName, speed, loc)
		require.NoError(t, err)
		assert.Equal(t, validCourierName, c.Name())
		assert.Equal(t, speed, c.Speed())
		assert.Equal(t, model.MustVolume(courierMaxVolume), c.MaxVolume())
		assert.Equal(t, loc, c.Location())
		assert.NotEqual(t, uuid.Nil, c.ID())
	})

	t.Run("empty name returns ErrValueIsRequired", func(t *testing.T) {
		_, err := NewCourier("", speed, loc)
		assert.ErrorIs(t, err, errs.ErrValueIsRequired)
	})

	t.Run("zero speed returns ErrValueIsRequired", func(t *testing.T) {
		_, err := NewCourier(validCourierName, Speed{}, loc)
		assert.ErrorIs(t, err, errs.ErrValueIsRequired)
	})

	t.Run("zero location returns ErrValueIsRequired", func(t *testing.T) {
		_, err := NewCourier(validCourierName, speed, model.Location{})
		assert.ErrorIs(t, err, errs.ErrValueIsRequired)
	})
}

func TestCourier_Move(t *testing.T) {
	c := validCourier(t)

	newLoc, err := model.NewLocation(3, 7)
	require.NoError(t, err)

	c.Move(newLoc)
	assert.Equal(t, newLoc, c.Location())
}

func TestCourier_CanTakeOrder(t *testing.T) {
	c := validCourier(t) // maxVolume = 20

	t.Run("returns true when volume fits", func(t *testing.T) {
		vol, err := model.NewVolume(15)
		require.NoError(t, err)
		assert.True(t, c.CanTakeOrder(vol))
	})

	t.Run("returns false when volume exceeds max", func(t *testing.T) {
		vol, err := model.NewVolume(21)
		require.NoError(t, err)
		assert.False(t, c.CanTakeOrder(vol))
	})

	t.Run("accounts for active assignments", func(t *testing.T) {
		assigned, err := model.NewVolume(15)
		require.NoError(t, err)
		require.NoError(t, c.AssignOrder(uuid.New(), assigned, validLocation(t)))

		// 15 assigned, max 20 — only 5 left
		six, err := model.NewVolume(6)
		require.NoError(t, err)
		assert.False(t, c.CanTakeOrder(six))

		five, err := model.NewVolume(5)
		require.NoError(t, err)
		assert.True(t, c.CanTakeOrder(five))
	})
}

func TestCourier_AssignOrder(t *testing.T) {
	t.Run("assigns order successfully", func(t *testing.T) {
		c := validCourier(t)
		vol, err := model.NewVolume(5)
		require.NoError(t, err)

		err = c.AssignOrder(uuid.New(), vol, validLocation(t))
		require.NoError(t, err)
		assert.Len(t, c.Assignments(), 1)
	})

	t.Run("returns ErrCannotTakeOrder when volume exceeded", func(t *testing.T) {
		c := validCourier(t)
		vol, err := model.NewVolume(21)
		require.NoError(t, err)

		err = c.AssignOrder(uuid.New(), vol, validLocation(t))
		assert.ErrorIs(t, err, ErrCannotTakeOrder)
	})
}

func TestCourier_CompleteAssignment(t *testing.T) {
	orderLoc, err := model.NewLocation(5, 5)
	require.NoError(t, err)
	orderVol, err := model.NewVolume(3)
	require.NoError(t, err)

	t.Run("completes assignment when courier is at same location", func(t *testing.T) {
		c, err := NewCourier(validCourierName, validSpeed(t), orderLoc)
		require.NoError(t, err)
		require.NoError(t, c.AssignOrder(uuid.New(), orderVol, orderLoc))

		assignmentID := c.Assignments()[0].ID()
		assert.NoError(t, c.CompleteAssignment(assignmentID))
		assert.True(t, c.Assignments()[0].Status().Equals(statusCompleted))
	})

	t.Run("completes assignment when courier is 1 cell away", func(t *testing.T) {
		nearLoc, err := model.NewLocation(5, 6)
		require.NoError(t, err)
		c, err := NewCourier(validCourierName, validSpeed(t), nearLoc)
		require.NoError(t, err)
		require.NoError(t, c.AssignOrder(uuid.New(), orderVol, orderLoc))

		assignmentID := c.Assignments()[0].ID()
		assert.NoError(t, c.CompleteAssignment(assignmentID))
	})

	t.Run("returns ErrCourierTooFar when courier is far", func(t *testing.T) {
		farLoc, err := model.NewLocation(1, 1)
		require.NoError(t, err)
		c, err := NewCourier(validCourierName, validSpeed(t), farLoc)
		require.NoError(t, err)
		require.NoError(t, c.AssignOrder(uuid.New(), orderVol, orderLoc))

		assignmentID := c.Assignments()[0].ID()
		assert.ErrorIs(t, c.CompleteAssignment(assignmentID), ErrCourierTooFar)
	})

	t.Run("returns ErrAssignmentNotFound for unknown id", func(t *testing.T) {
		c := validCourier(t)
		assert.ErrorIs(t, c.CompleteAssignment(uuid.New()), ErrAssignmentNotFound)
	})
}
