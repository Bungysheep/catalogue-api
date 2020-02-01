package catalogue

import "time"

// Catalogue type
type Catalogue struct {
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Details     string    `json:"details"`
	Status      string    `json:"status"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedBy  string    `json:"modified_by"`
	ModifiedAt  time.Time `json:"modified_at"`
	Vers        int64     `json:"vers"`
}

// NewCatalogue - Creates catalogue
func NewCatalogue() *Catalogue {
	return &Catalogue{}
}

// GetCode - Returns catalogue code
func (clg *Catalogue) GetCode() string {
	return clg.Code
}

// GetDescription - Returns catalogue description
func (clg *Catalogue) GetDescription() string {
	return clg.Description
}

// GetDetails - Returns catalogue details
func (clg *Catalogue) GetDetails() string {
	return clg.Details
}

// GetStatus - Returns catalogue status
func (clg *Catalogue) GetStatus() string {
	return clg.Status
}

// GetCreatedBy - Returns created by
func (clg *Catalogue) GetCreatedBy() string {
	return clg.CreatedBy
}

// GetCreatedAt - Returns created at
func (clg *Catalogue) GetCreatedAt() time.Time {
	return clg.CreatedAt
}

// GetModifiedBy - Returns modified by
func (clg *Catalogue) GetModifiedBy() string {
	return clg.ModifiedBy
}

// GetModifiedAt - Returns modified at
func (clg *Catalogue) GetModifiedAt() time.Time {
	return clg.ModifiedAt
}

// GetVers - Returns vers
func (clg *Catalogue) GetVers() int64 {
	return clg.Vers
}

// DoValidate - Validate catalogue
func (clg *Catalogue) DoValidate() (bool, string) {
	return true, ""
}
