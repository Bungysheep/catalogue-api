package product

import (
	"time"

	"github.com/bungysheep/catalogue-api/pkg/models/v1/unitofmeasure"
)

// Product type
type Product struct {
	ID             int64                          `json:"id"`
	CatalogueCode  string                         `json:"clg_code"`
	Code           string                         `json:"code"`
	Description    string                         `json:"description"`
	Details        string                         `json:"details"`
	Status         string                         `json:"status"`
	CreatedBy      string                         `json:"created_by"`
	CreatedAt      time.Time                      `json:"created_at"`
	ModifiedBy     string                         `json:"modified_by"`
	ModifiedAt     time.Time                      `json:"modified_at"`
	Vers           int64                          `json:"vers"`
	UnitOfMeasures []*unitofmeasure.UnitOfMeasure `json:"uoms"`
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
	return prod.Code
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
	return prod.Status
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
		if uom.GetIsDefault() {
			return uom
		}
	}

	return nil
}

// GetUoms - Returns product uoms
func (prod *Product) GetUoms() []*unitofmeasure.UnitOfMeasure {
	return prod.UnitOfMeasures
}
