package changemode

// ChangeMode type
type ChangeMode int

const (
	// Unchange change mode
	Unchange ChangeMode = iota

	// Add change mode
	Add

	// Update change mode
	Update

	// Delete change mode
	Delete
)

func (cm ChangeMode) String() string {
	return [...]string{"A", "I"}[cm]
}
