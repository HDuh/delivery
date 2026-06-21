package courier

import (
	"errors"

	"github.com/google/uuid"

	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/ddd"
)

var (
	ErrVolumeIsZero               = errors.New("volume should not be 0")
	ErrInvalidOrderID             = errors.New("orderID should not be empty")
	ErrAssignmentAlreadyCompleted = errors.New("assignment already completed")
	ErrCourierTooFar              = errors.New("courier is too far from order location")
)

type Assignment struct {
	baseEntity ddd.BaseEntity[uuid.UUID]
	orderId    uuid.UUID
	volume     model.Volume
	location   model.Location
	status     AssignmentStatus
}

func NewAssignment(orderID uuid.UUID, volume model.Volume, location model.Location) (*Assignment, error) {
	if orderID == uuid.Nil {
		return nil, ErrInvalidOrderID
	}
	if volume.IsZero() {
		return nil, ErrVolumeIsZero
	}

	return &Assignment{
		baseEntity: *ddd.NewBaseEntity(uuid.New()),
		orderId:    orderID,
		volume:     volume,
		location:   location,
		status:     statusAssigned,
	}, nil
}

func (a *Assignment) ID() uuid.UUID {
	return a.baseEntity.ID()
}

func (a *Assignment) OrderID() uuid.UUID {
	return a.orderId
}

func (a *Assignment) Volume() model.Volume {
	return a.volume
}

func (a *Assignment) Location() model.Location {
	return a.location
}

func (a *Assignment) Status() AssignmentStatus {
	return a.status
}

func (a *Assignment) Complete(courierLocation model.Location) error {
	if a.status.Equals(statusCompleted) {
		return ErrAssignmentAlreadyCompleted
	}
	if courierLocation.DistanceTo(a.location) > 1 {
		return ErrCourierTooFar
	}
	a.status = statusCompleted
	return nil
}
