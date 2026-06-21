package courier

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSpeed(t *testing.T) {
	t.Run("returns ErrInvalidSpeed when speed is 0", func(t *testing.T) {
		_, err := NewSpeed(0)
		assert.ErrorIs(t, err, ErrInvalidSpeed)
	})

	t.Run("returns ErrInvalidSpeed when speed exceeds max", func(t *testing.T) {
		_, err := NewSpeed(maxSpeed + 1)
		assert.ErrorIs(t, err, ErrInvalidSpeed)
	})

	for _, v := range []uint8{minSpeed, 2, 3, maxSpeed} {
		v := v
		t.Run("valid speed returns Speed", func(t *testing.T) {
			s, err := NewSpeed(v)
			require.NoError(t, err)
			assert.Equal(t, v, s.Value())
		})
	}
}

func TestSpeed_Equal(t *testing.T) {
	s1, err := NewSpeed(3)
	require.NoError(t, err)

	t.Run("equal speeds", func(t *testing.T) {
		s2, err := NewSpeed(3)
		require.NoError(t, err)
		assert.True(t, s1.Equal(s2))
	})

	t.Run("different speeds", func(t *testing.T) {
		s2, err := NewSpeed(1)
		require.NoError(t, err)
		assert.False(t, s1.Equal(s2))
	})
}
