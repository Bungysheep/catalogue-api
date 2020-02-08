package product

import (
	"strings"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/models/v1/basemodel"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/productcustomfield"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/unitofmeasure"
)

// Product type
type Product struct {
	basemodel.BaseModel
	ID             int64                                    `json:"id"`
	CatalogueCode  string                                   `json:"clg_code"`
	Code           string                                   `json:"code"`
	Description    string                                   `json:"description" mandatory:"true" max_length:"32"`
	Details        string                                   `json:"details" max_length:"64"`
	Status         string                                   `json:"status" mandatory:"true" max_length:"1"`
	CreatedBy      string                                   `json:"created_by"`
	CreatedAt      time.Time                                `json:"created_at"`
	ModifiedBy     string                                   `json:"modified_by"`
	ModifiedAt     time.Time                                `json:"modified_at"`
	Vers           int64                                    `json:"vers"`
	UnitOfMeasures []*unitofmeasure.UnitOfMeasure           `json:"uoms"`
	CustomFields   []*productcustomfield.ProductCustomField `json:"custom_fields"`
}

// NewProduct - Creates product
func NewProduct() *Product {
	return &Product{}
}

// GetID - Returns product id
func (prod *Product) GetID() int64 {
	return prod.ID
}

// GetCatalogueCode - Returns catalogue code
func (prod *Product) GetCatalogueCode() string {
	return prod.CatalogueCode
}

// GetCode - Returns product code
func (prod *Product) GetCode() string {
	return strings.ToUpper(prod.Code)
}

// GetDescription - Returns product description
func (prod *Product) GetDescription() string {
	return prod.Description
}

// GetDetails - Returns product details
func (prod *Product) GetDetails() string {
	return prod.Details
}

// GetStatus - Returns product status
func (prod *Product) GetStatus() string {
	return strings.ToUpper(prod.Status)
}

// GetCreatedBy - Returns created by
func (prod *Product) GetCreatedBy() string {
	return prod.CreatedBy
}

// GetCreatedAt - Returns created at
func (prod *Product) GetCreatedAt() time.Time {
	return prod.CreatedAt
}

// GetModifiedBy - Returns modified by
func (prod *Product) GetModifiedBy() string {
	return prod.ModifiedBy
}

// GetModifiedAt - Returns modified at
func (prod *Product) GetModifiedAt() time.Time {
	return prod.ModifiedAt
}

// GetVers - Returns vers
func (prod *Product) GetVers() int64 {
	return prod.Vers
}

// GetDefaultUom - Returns default uom
func (prod *Product) GetDefaultUom() *unitofmeasure.UnitOfMeasure {
	for _, uom := range prod.UnitOfMeasures {
		if uom.IsDefault() {
			return uom
		}
	}

	return nil
}

// GetAllUoms - Returns all product uoms
func (prod *Product) GetAllUoms() []*unitofmeasure.UnitOfMeasure {
	return prod.UnitOfMeasures
}

// GetAllCustomFields - Returns all product custom fields
func (prod *Product) GetAllCustomFields() []*productcustomfield.ProductCustomField {
	return prod.CustomFields
}

// GetUom - Returns uom
func (prod *Product) GetUom(uomID int64) *unitofmeasure.UnitOfMeasure {
	for _, uom := range prod.UnitOfMeasures {
		if uom.GetID() == uomID {
			return uom
		}
	}

	return nil
}

// GetCustomField - Returns custom field
func (prod *Product) GetCustomField(fieldID int64) *productcustomfield.ProductCustomField {
	for _, field := range prod.CustomFields {
		if field.GetID() == fieldID {
			return field
		}
	}

	return nil
}

// GetNumberOfDefaultUom - Returns number of default uom
func (prod *Product) GetNumberOfDefaultUom() int {
	count := 0

	for _, uom := range prod.UnitOfMeasures {
		if uom.IsDefault() {
			count++
		}

		if count > 1 {
			break
		}
	}

	return count
}

// DoValidate - Validate product
func (prod *Product) DoValidate(otherProd *Product) (bool, string) {
	var ok bool
	var message string

	ok, message = prod.DoValidateBase(*prod)
	if !ok {
		return false, message
	}

	nbrDefaultUom := prod.GetNumberOfDefaultUom()
	if nbrDefaultUom == 0 {
		// If existing default uom has been changed to be non-default uom
		if otherProd == nil || prod.GetUom(otherProd.GetDefaultUom().GetID()) != nil {
			return false, "No default unit of measure."
		}
	} else if nbrDefaultUom == 1 {
		// If existing default uom is different with updated default uom
		if otherProd != nil && prod.GetDefaultUom().GetID() != otherProd.GetDefaultUom().GetID() {
			return false, "Found multiple default unit of measure."
		}
	} else if nbrDefaultUom > 1 {
		return false, "Found multiple default unit of measure."
	}

	for _, uom := range prod.UnitOfMeasures {
		ok, message = uom.DoValidate()
		if !ok {
			return false, message
		}
	}

	return true, ""
}
