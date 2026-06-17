package model

import (
	"delivery/internal/pkg/errs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_VolumeBeCorrectWhenParamsAreCorrectOnCreated(t *testing.T) {
	// Arrange

	// Act
	volume, err := NewVolume(1)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, volume)
	assert.Equal(t, 1, volume.Value())
}

func Test_VolumeReturnErrorWhenValueIsNegative(t *testing.T) {
	// Act
	_, err := NewVolume(-1)

	// Assert
	expected := errs.NewValueMustBeGreaterThan("value", -1, minVolumeValue)

	if err.Error() != expected.Error() {
		t.Errorf("expected %v, got %v", expected, err)
	}
}

func Test_VolumeAddShouldBeSuccessWhenParamsAreCorrect(t *testing.T) {
	// Arrange
	v1 := MustVolume(5)
	v2 := MustVolume(3)

	// Act
	result := v1.Add(v2)

	// Assert
	assert.Equal(t, 8, result.Value())
}

func Test_VolumeSubtractShouldBeSuccessWhenEnoughValue(t *testing.T) {
	// Arrange
	v1 := MustVolume(5)
	v2 := MustVolume(3)

	// Act
	result, err := v1.Subtract(v2)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 2, result.Value())
}

func Test_QuantitySubtractReturnErrorWhenNotEnoughValue(t *testing.T) {
	// Arrange
	v1 := MustVolume(2)
	v2 := MustVolume(5)

	// Act
	_, err := v1.Subtract(v2)

	// Assert
	assert.Error(t, err)
}
