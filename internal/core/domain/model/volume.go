package model

import (
	"delivery/internal/pkg/errs"
)

const minVolumeValue = 0

var QuantityZero = MustVolume(minVolumeValue)

type Volume struct {
	value int
}

func NewVolume(value int) (Volume, error) {
	if value < minVolumeValue {
		return Volume{}, errs.NewValueMustBeGreaterThan("value", value, minVolumeValue)
	}

	return Volume{
		value: value,
	}, nil
}

func MustVolume(value int) Volume {
	v, err := NewVolume(value)
	if err != nil {
		panic(err)
	}
	return v
}

func (v Volume) Value() int {
	return v.value
}

func (v Volume) Add(other Volume) Volume {
	return MustVolume(v.value + other.value)
}

func (v Volume) Subtract(other Volume) (Volume, error) {
	if other.value > v.value {
		return Volume{}, errs.NewValueMustBeLessOrEqual(
			"quantity",
			other.value,
			v.value,
		)
	}

	return NewVolume(v.value - other.value)
}

func (v Volume) Equal(other Volume) bool {
	return v.value == other.value
}

func (v Volume) IsZero() bool {
	return v.value == 0
}

func (v Volume) GreaterThan(other Volume) bool {
	return v.value > other.value
}

func (v Volume) GreaterThanOrEqual(other Volume) bool {
	return v.value >= other.value
}

func (v Volume) LessThan(other Volume) bool {
	return v.value < other.value
}

func (v Volume) LessThanOrEqual(other Volume) bool {
	return v.value <= other.value
}
