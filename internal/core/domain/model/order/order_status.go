package order

const (
	StatusEmpty     OrderStatus = ""
	StatusCreated   OrderStatus = "Created"
	statusCompleted OrderStatus = "Completed"
	statusAssigned  OrderStatus = "Assigned"
)

type OrderStatus string

func (s OrderStatus) Equal(other OrderStatus) bool {
	return s == other
}

func (s OrderStatus) IsValid() bool {
	return s != StatusEmpty
}

func (s OrderStatus) String() string {
	return string(s)
}
