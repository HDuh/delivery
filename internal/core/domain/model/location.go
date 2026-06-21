package model

import (
	"delivery/internal/pkg/errs"
	"fmt"
	"math"
)

const (
	minCoordinate = 1
	maxCoordinate = 10
)

type Location struct {
	xCoord uint8
	yCoord uint8
}

func (l Location) XCoord() uint8 {
	return l.xCoord
}

func (l Location) YCoord() uint8 {
	return l.yCoord
}

func NewLocation(x uint8, y uint8) (Location, error) {
	if x < minCoordinate || x > maxCoordinate {
		return Location{}, errs.NewMustBeBetween("location.x", x, minCoordinate, maxCoordinate)
	}
	if y < minCoordinate || y > maxCoordinate {
		return Location{}, errs.NewMustBeBetween("location.y", y, minCoordinate, maxCoordinate)
	}

	return Location{
		xCoord: x,
		yCoord: y,
	}, nil
}

func (l Location) IsZero() bool { return l.xCoord == 0 && l.yCoord == 0 }

func (l Location) Equal(other Location) bool {
	return l.xCoord == other.xCoord && l.yCoord == other.yCoord
}

func (l Location) DistanceTo(other Location) uint8 {
	return (uint8(math.Abs(float64(int(l.xCoord)-int(other.xCoord)))) +
		uint8((math.Abs(float64(int(l.yCoord) - int(other.yCoord))))))
}

func (l Location) String() string {
	return fmt.Sprintf("(%d, %d)", l.xCoord, l.yCoord)
}
