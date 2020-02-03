package definitiontype

// DefinitionType type
type DefinitionType int

const (
	// Alphanumeric definition type
	Alphanumeric DefinitionType = iota

	// Numeric definition type
	Numeric

	// Date definition type
	Date
)

func (dt DefinitionType) String() string {
	return [...]string{"A", "N", "D"}[dt]
}
