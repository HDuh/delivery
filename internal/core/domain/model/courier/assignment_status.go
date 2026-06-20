package courier

import "fmt"

type AssignmentStatus struct {
	name string
}

var (
	statusEmpty     = AssignmentStatus{}
	statusAssigned  = AssignmentStatus{name: "assigned"}
	statusCompleted = AssignmentStatus{name: "completed"}
)

func AssignmentStatusAssigned() AssignmentStatus  { return statusAssigned }
func AssignmentStatusCompleted() AssignmentStatus { return statusCompleted }

func (s AssignmentStatus) IsEmpty() bool                      { return s == statusEmpty }
func (s AssignmentStatus) Equals(other AssignmentStatus) bool { return s == other }
func (s AssignmentStatus) String() string                     { return s.name }

func ParseAssignmentStatus(v string) (AssignmentStatus, error) {
	switch v {
	case statusAssigned.name:
		return statusAssigned, nil
	case statusCompleted.name:
		return statusCompleted, nil
	default:
		return statusEmpty, fmt.Errorf("unknown assignment status: %q", v)
	}
}
