package status

// Status type
type Status int

const (
	// Active status
	Active Status = iota

	// Inactive status
	Inactive
)

func (s Status) String() string {
	return [...]string{"A", "I"}[s]
}
