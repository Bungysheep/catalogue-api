package unitofmeasure

// UnitOfMeasure type
type UnitOfMeasure struct {
	ID          int64   `json:"id"`
	ProdID      int64   `json:"prod_id"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
	IsDefault   bool    `json:"is_default"`
	Ratio       float64 `json:"ratio"`
	Vers        int64   `json:"vers"`
}

// NewUnitOfMeasure - Creates unit of measure
func NewUnitOfMeasure() *UnitOfMeasure {
	return &UnitOfMeasure{}
}

// GetID - Returns uom id
func (uom *UnitOfMeasure) GetID() int64 {
	return uom.ID
}

// GetProdID - Returns prod id
func (uom *UnitOfMeasure) GetProdID() int64 {
	return uom.ProdID
}

// GetCode - Returns uom code
func (uom *UnitOfMeasure) GetCode() string {
	return uom.Code
}

// GetDescription - Returns uom description
func (uom *UnitOfMeasure) GetDescription() string {
	return uom.Description
}

// GetIsDefault - Returns whether it is default uom or not
func (uom *UnitOfMeasure) GetIsDefault() bool {
	return uom.IsDefault
}

// GetRatio - Returns uom ratio
func (uom *UnitOfMeasure) GetRatio() float64 {
	return uom.Ratio
}

// GetVers - Returns vers
func (uom *UnitOfMeasure) GetVers() int64 {
	return uom.Vers
}
