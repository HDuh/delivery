package model

import (
	"delivery/internal/pkg/errs"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLocation(t *testing.T) {
	t.Run("valid coordinates", func(t *testing.T) {
		loc, err := NewLocation(1, 1)
		assert.NoError(t, err)
		assert.Equal(t, uint8(1), loc.XCoord())
		assert.Equal(t, uint8(1), loc.YCoord())
	})

	t.Run("valid boundary max", func(t *testing.T) {
		loc, err := NewLocation(10, 10)
		assert.NoError(t, err)
		assert.Equal(t, uint8(10), loc.XCoord())
		assert.Equal(t, uint8(10), loc.YCoord())
	})

	t.Run("x is zero", func(t *testing.T) {
		_, err := NewLocation(0, 5)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrMustBeBetween))
	})

	t.Run("y is zero", func(t *testing.T) {
		_, err := NewLocation(5, 0)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrMustBeBetween))
	})

	t.Run("x exceeds max", func(t *testing.T) {
		_, err := NewLocation(11, 5)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrMustBeBetween))
	})

	t.Run("y exceeds max", func(t *testing.T) {
		_, err := NewLocation(5, 11)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errs.ErrMustBeBetween))
	})
}

func TestLocation_Equal(t *testing.T) {
	t.Run("same coordinates", func(t *testing.T) {
		loc1, _ := NewLocation(3, 5)
		loc2, _ := NewLocation(3, 5)
		assert.True(t, loc1.Equal(loc2))
	})

	t.Run("different x", func(t *testing.T) {
		loc1, _ := NewLocation(3, 5)
		loc2, _ := NewLocation(4, 5)
		assert.False(t, loc1.Equal(loc2))
	})

	t.Run("different y", func(t *testing.T) {
		loc1, _ := NewLocation(3, 5)
		loc2, _ := NewLocation(3, 6)
		assert.False(t, loc1.Equal(loc2))
	})
}

func TestLocation_DistanceTo(t *testing.T) {
	t.Run("same location", func(t *testing.T) {
		loc, _ := NewLocation(5, 5)
		assert.Equal(t, uint8(0), loc.DistanceTo(loc))
	})

	t.Run("differs only in x", func(t *testing.T) {
		loc1, _ := NewLocation(2, 5)
		loc2, _ := NewLocation(5, 5)
		assert.Equal(t, uint8(3), loc1.DistanceTo(loc2))
	})

	t.Run("differs only in y", func(t *testing.T) {
		loc1, _ := NewLocation(5, 2)
		loc2, _ := NewLocation(5, 7)
		assert.Equal(t, uint8(5), loc1.DistanceTo(loc2))
	})

	t.Run("differs in both axes", func(t *testing.T) {
		loc1, _ := NewLocation(2, 6)
		loc2, _ := NewLocation(4, 9)
		assert.Equal(t, uint8(5), loc1.DistanceTo(loc2))
	})

	t.Run("symmetric", func(t *testing.T) {
		loc1, _ := NewLocation(2, 6)
		loc2, _ := NewLocation(4, 9)
		assert.Equal(t, loc1.DistanceTo(loc2), loc2.DistanceTo(loc1))
	})

	t.Run("max distance", func(t *testing.T) {
		loc1, _ := NewLocation(1, 1)
		loc2, _ := NewLocation(10, 10)
		assert.Equal(t, uint8(18), loc1.DistanceTo(loc2))
	})
}
