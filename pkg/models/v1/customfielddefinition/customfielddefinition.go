package customfielddefinition

import (
	"time"

	"github.com/bungysheep/catalogue-api/pkg/models/v1/basemodel"
)

// CustomFieldDefinition type
type CustomFieldDefinition struct {
	basemodel.BaseModel
	ID            int64     `json:"id"`
	CatalogueCode string    `json:"clg_code"`
	Caption       string    `json:"caption" mandatory:"true" max_length:"32"`
	Type          string    `json:"type" mandatory:"true" max_length:"1" valid_value:"A,N,D"`
	CreatedBy     string    `json:"created_by"`
	CreatedAt     time.Time `json:"created_at"`
	ModifiedBy    string    `json:"modified_by"`
	ModifiedAt    time.Time `json:"modified_at"`
	Vers          int64     `json:"vers"`
}

// NewCustomFieldDefinition - Creates custom field definition
func NewCustomFieldDefinition() *CustomFieldDefinition {
	return &CustomFieldDefinition{}
}

// GetID - Returns custom field definition id
func (cfd *CustomFieldDefinition) GetID() int64 {
	return cfd.ID
}

// GetCatalogueCode - Returns catalogue code
func (cfd *CustomFieldDefinition) GetCatalogueCode() string {
	return cfd.CatalogueCode
}

// GetCaption - Returns caption
func (cfd *CustomFieldDefinition) GetCaption() string {
	return cfd.Caption
}

// GetType - Returns type
func (cfd *CustomFieldDefinition) GetType() string {
	return cfd.Type
}

// GetCreatedBy - Returns created by
func (cfd *CustomFieldDefinition) GetCreatedBy() string {
	return cfd.CreatedBy
}

// GetCreatedAt - Returns created at
func (cfd *CustomFieldDefinition) GetCreatedAt() time.Time {
	return cfd.CreatedAt
}

// GetModifiedBy - Returns modified by
func (cfd *CustomFieldDefinition) GetModifiedBy() string {
	return cfd.ModifiedBy
}

// GetModifiedAt - Returns modified at
func (cfd *CustomFieldDefinition) GetModifiedAt() time.Time {
	return cfd.ModifiedAt
}

// GetVers - Returns vers
func (cfd *CustomFieldDefinition) GetVers() int64 {
	return cfd.Vers
}

// IsEqual - Whether equal
func (cfd *CustomFieldDefinition) IsEqual(otherFieldDef *CustomFieldDefinition) bool {
	return cfd.ID == otherFieldDef.GetID() &&
		cfd.Caption == otherFieldDef.GetCaption() &&
		cfd.Type == otherFieldDef.GetType()
}

// DoValidate - Validate custom field definition
func (cfd *CustomFieldDefinition) DoValidate() (bool, string) {
	return cfd.DoValidateBase(*cfd)
}
