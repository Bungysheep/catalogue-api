package productcustomfield

import (
	"time"

	"github.com/bungysheep/catalogue-api/pkg/models/v1/basemodel"
)

// ProductCustomField type
type ProductCustomField struct {
	basemodel.BaseModel
	ID           int64     `json:"id"`
	ProdID       int64     `json:"prod_id"`
	FieldID      int64     `json:"field_id"`
	AlphaValue   string    `json:"alpha_value" max_length:"16"`
	NumericValue float64   `json:"numeric_value"`
	DateValue    time.Time `json:"date_value"`
}

// NewProductCustomField - Creates product custom field
func NewProductCustomField() *ProductCustomField {
	return &ProductCustomField{}
}

// GetID - Returns product custom field id
func (pcf *ProductCustomField) GetID() int64 {
	return pcf.ID
}

// GetProdID - Returns prod id
func (pcf *ProductCustomField) GetProdID() int64 {
	return pcf.ProdID
}

// GetFieldID - Returns field id
func (pcf *ProductCustomField) GetFieldID() int64 {
	return pcf.FieldID
}

// GetAlphaValue - Returns alpha value
func (pcf *ProductCustomField) GetAlphaValue() string {
	return pcf.AlphaValue
}

// GetNumericValue - Returns numeric value
func (pcf *ProductCustomField) GetNumericValue() float64 {
	return pcf.NumericValue
}

// GetDateValue - Returns date value
func (pcf *ProductCustomField) GetDateValue() time.Time {
	return pcf.DateValue
}
