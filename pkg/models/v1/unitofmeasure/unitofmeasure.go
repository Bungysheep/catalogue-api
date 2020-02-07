package unitofmeasure

import (
	"strings"

	"github.com/bungysheep/catalogue-api/pkg/commons/changemode"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/basemodel"
)

// UnitOfMeasure type
type UnitOfMeasure struct {
	basemodel.BaseModel
	ID          int64                 `json:"id"`
	ProdID      int64                 `json:"prod_id"`
	Code        string                `json:"code" mandatory:"true" max_length:"16"`
	Description string                `json:"description" mandatory:"true" max_length:"32"`
	Ratio       float64               `json:"ratio"`
	Vers        int64                 `json:"vers"`
	ChangeMode  changemode.ChangeMode `json:"change_mode"`
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
	return strings.ToUpper(uom.Code)
}

// GetDescription - Returns uom description
func (uom *UnitOfMeasure) GetDescription() string {
	return uom.Description
}

// GetRatio - Returns uom ratio
func (uom *UnitOfMeasure) GetRatio() float64 {
	return uom.Ratio
}

// GetVers - Returns vers
func (uom *UnitOfMeasure) GetVers() int64 {
	return uom.Vers
}

// GetChangeMode - Returns change mode
func (uom *UnitOfMeasure) GetChangeMode() changemode.ChangeMode {
	return uom.ChangeMode
}

// IsDefault - Whether default uom or not
func (uom *UnitOfMeasure) IsDefault() bool {
	return uom.Ratio == 1
}

// IsEqual - Whether equal
func (uom *UnitOfMeasure) IsEqual(otherUom *UnitOfMeasure) bool {
	return uom.ID == otherUom.GetID() &&
		uom.Code == otherUom.GetCode() &&
		uom.Description == otherUom.GetDescription() &&
		uom.Ratio == otherUom.GetRatio()
}

// DoValidate - Validate uom
func (uom *UnitOfMeasure) DoValidate() (bool, string) {
	return uom.DoValidateBase(*uom)
}
