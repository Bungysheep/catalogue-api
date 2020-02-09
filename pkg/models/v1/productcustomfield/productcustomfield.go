package productcustomfield

import (
	"fmt"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/commons/definitiontype"
	"github.com/bungysheep/catalogue-api/pkg/configs"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/basemodel"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/customfielddefinition"
)

// ProductCustomField type
type ProductCustomField struct {
	basemodel.BaseModel
	ID           int64     `json:"id"`
	ProdID       int64     `json:"prod_id"`
	FieldID      int64     `json:"field_id"`
	AlphaValue   string    `json:"alpha_value" max_length:"64"`
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

// DoValidate - Validate product custom field
func (pcf *ProductCustomField) DoValidate(fieldDef *customfielddefinition.CustomFieldDefinition) (bool, string) {
	defaultDate, _ := time.Parse(configs.DATEFORMAT, configs.DEFAULTDATE)

	if fieldDef == nil {
		return false, fmt.Sprintf("Custom Field '%d' has invalid definition.", pcf.GetFieldID())
	}

	switch fieldDef.GetType() {
	case definitiontype.Alphanumeric.String():
		pcf.NumericValue = 0
		pcf.DateValue = defaultDate

		if fieldDef.GetMandatory() && pcf.GetAlphaValue() == "" {
			return false, fmt.Sprintf("Custom Field '%d' must be specified.", pcf.GetFieldID())
		}

	case definitiontype.Numeric.String():
		pcf.AlphaValue = ""
		pcf.DateValue = defaultDate

		if fieldDef.GetMandatory() && pcf.GetNumericValue() == 0 {
			return false, fmt.Sprintf("Custom Field '%d' must be specified.", pcf.GetFieldID())
		}

	case definitiontype.Date.String():
		pcf.AlphaValue = ""
		pcf.NumericValue = 0

		if fieldDef.GetMandatory() && pcf.GetDateValue() == defaultDate {
			return false, fmt.Sprintf("Custom Field '%d' must be specified.", pcf.GetFieldID())
		}

	}

	return pcf.DoValidateBase(*pcf)
}
