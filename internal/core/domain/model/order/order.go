package order

import (
	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/errs"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrCannotAssign   = errors.New("order cannot be assigned")
	ErrCannotComplete = errors.New("order cannot be completed")
)

type Order struct {
	id       uuid.UUID
	location model.Location
	volume   model.Volume
	status   OrderStatus
}

func NewOrder(id uuid.UUID, location model.Location, volume model.Volume) (*Order, error) {
	if id == uuid.Nil {
		return nil, errs.NewValueIsRequired("id")
	}
	return &Order{
		id:       id,
		location: location,
		volume:   volume,
		status:   StatusCreated,
	}, nil
}

func (o *Order) ID() uuid.UUID            { return o.id }
func (o *Order) Location() model.Location { return o.location }
func (o *Order) Volume() model.Volume     { return o.volume }
func (o *Order) Status() OrderStatus      { return o.status }

func (o *Order) Assign() error {
	if o.status != StatusCreated {
		return fmt.Errorf("%w: current status is %s, expected %s", ErrCannotAssign, o.status, StatusCreated)
	}
	o.status = statusAssigned
	return nil
}

func (o *Order) Complete() error {
	if o.status != statusAssigned {
		return fmt.Errorf("%w: current status is %s, expected %s", ErrCannotComplete, o.status, statusAssigned)
	}
	o.status = statusCompleted
	return nil
}
