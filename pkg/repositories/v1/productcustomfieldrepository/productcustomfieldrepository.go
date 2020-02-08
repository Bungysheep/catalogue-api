package productcustomfieldrepository

import (
	"context"

	productcustomfieldmodel "github.com/bungysheep/catalogue-api/pkg/models/v1/productcustomfield"
)

// IProductCustomFieldRepository type
type IProductCustomFieldRepository interface {
	GetByID(context.Context, int64) (*productcustomfieldmodel.ProductCustomField, error)
	GetByProduct(context.Context, int64) ([]*productcustomfieldmodel.ProductCustomField, error)
	Create(context.Context, *productcustomfieldmodel.ProductCustomField) (int64, error)
	Update(context.Context, *productcustomfieldmodel.ProductCustomField) (int64, error)
	Delete(context.Context, int64) (int64, error)
	DeleteByProduct(context.Context, int64) error
}

type productCustomFieldRepository struct {
}

// NewProductCustomFieldRepository - Create product custom field repository
func NewProductCustomFieldRepository() IProductCustomFieldRepository {
	return &productCustomFieldRepository{}
}

func (pcfRepo *productCustomFieldRepository) GetByID(ctx context.Context, id int64) (*productcustomfieldmodel.ProductCustomField, error) {
	return nil, nil
}

func (pcfRepo *productCustomFieldRepository) GetByProduct(ctx context.Context, id int64) ([]*productcustomfieldmodel.ProductCustomField, error) {
	return nil, nil
}

func (pcfRepo *productCustomFieldRepository) Create(ctx context.Context, data *productcustomfieldmodel.ProductCustomField) (int64, error) {
	return 0, nil
}

func (pcfRepo *productCustomFieldRepository) Update(ctx context.Context, data *productcustomfieldmodel.ProductCustomField) (int64, error) {
	return 0, nil
}

func (pcfRepo *productCustomFieldRepository) Delete(ctx context.Context, id int64) (int64, error) {
	return 0, nil
}

func (pcfRepo *productCustomFieldRepository) DeleteByProduct(ctx context.Context, prodID int64) error {
	return nil
}
