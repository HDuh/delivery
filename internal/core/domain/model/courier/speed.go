package courier

import "fmt"

const (
	minSpeed = 1
	maxSpeed = 4
)

var ErrInvalidSpeed = fmt.Errorf("speed must be between %d and %d", minSpeed, maxSpeed)

type Speed struct {
	value uint8
}

func NewSpeed(value uint8) (Speed, error) {
	if value < minSpeed || value > maxSpeed {
		return Speed{}, ErrInvalidSpeed
	}
	return Speed{value: value}, nil
}

func (s Speed) Value() uint8 { return s.value }

func (s Speed) IsZero() bool           { return s.value == 0 }
func (s Speed) Equal(other Speed) bool { return s.value == other.value }
