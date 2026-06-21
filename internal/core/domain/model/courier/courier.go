package courier

import (
	"errors"

	"delivery/internal/core/domain/model"
	"delivery/internal/pkg/ddd"
	"delivery/internal/pkg/errs"

	"github.com/google/uuid"
)

const courierMaxVolume = 20

var (
	ErrAssignmentNotFound = errors.New("assignment not found")
	ErrCannotTakeOrder    = errors.New("volume limit exceeded")
)

type Courier struct {
	baseAggregate *ddd.BaseAggregate[uuid.UUID]
	name          string
	location      model.Location
	speed         Speed
	maxVolume     model.Volume
	assignments   []*Assignment
}

func NewCourier(name string, speed Speed, location model.Location) (*Courier, error) {
	if name == "" {
		return nil, errs.NewValueIsRequired("name")
	}
	if speed.IsZero() {
		return nil, errs.NewValueIsRequired("speed")
	}
	if location.IsZero() {
		return nil, errs.NewValueIsRequired("location")
	}
	return &Courier{
		baseAggregate: ddd.NewBaseAggregate(uuid.New()),
		name:          name,
		location:      location,
		speed:         speed,
		maxVolume:     model.MustVolume(courierMaxVolume),
		assignments:   make([]*Assignment, 0),
	}, nil
}

func (c *Courier) ID() uuid.UUID              { return c.baseAggregate.ID() }
func (c *Courier) Name() string               { return c.name }
func (c *Courier) Location() model.Location   { return c.location }
func (c *Courier) Speed() Speed               { return c.speed }
func (c *Courier) MaxVolume() model.Volume    { return c.maxVolume }
func (c *Courier) Assignments() []*Assignment { return c.assignments }

func (c *Courier) Move(location model.Location) {
	c.location = location
}

func (c *Courier) CanTakeOrder(volume model.Volume) bool {
	total := model.QuantityZero
	for _, a := range c.assignments {
		if !a.Status().Equals(statusCompleted) {
			total = total.Add(a.Volume())
		}
	}
	return total.Add(volume).LessThanOrEqual(c.maxVolume)
}

func (c *Courier) AssignOrder(orderID uuid.UUID, volume model.Volume, location model.Location) error {
	if !c.CanTakeOrder(volume) {
		return ErrCannotTakeOrder
	}
	a, err := NewAssignment(orderID, volume, location)
	if err != nil {
		return err
	}
	c.assignments = append(c.assignments, a)
	return nil
}

func (c *Courier) CompleteAssignment(assignmentID uuid.UUID) error {
	for _, a := range c.assignments {
		if a.ID() == assignmentID {
			return a.Complete(c.location)
		}
	}
	return ErrAssignmentNotFound
}
