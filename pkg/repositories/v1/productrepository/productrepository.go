package productrepository

import (
	"context"
	"fmt"

	productmodel "github.com/bungysheep/catalogue-api/pkg/models/v1/product"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/productcustomfieldrepository"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/unitofmeasurerepository"
)

// IProductRepository type
type IProductRepository interface {
	GetByID(context.Context, int64) (*productmodel.Product, error)
	GetByCatalogue(context.Context, string) ([]*productmodel.Product, error)
	Create(context.Context, *productmodel.Product) (int64, error)
	Update(context.Context, *productmodel.Product) (int64, error)
	Delete(context.Context, int64) (int64, error)
	DeleteByCatalogue(context.Context, string) error
}

type productRepository struct {
}

// NewProductRepository - Create product repository
func NewProductRepository() IProductRepository {
	return &productRepository{}
}

func (prodRepo *productRepository) GetByID(ctx context.Context, id int64) (*productmodel.Product, error) {
	result := productmodel.NewProduct()

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, clg_code, code, descr, details, status, created_by, created_at, modified_by, modified_at, vers
		FROM products 
		WHERE id=$1`)
	if err != nil {
		return nil, fmt.Errorf("Failed preparing read product, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Failed reading product, error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("Failed retrieve product record, error: %v", err)
		}
		return nil, nil
	}

	if err := rows.Scan(
		&result.ID,
		&result.CatalogueCode,
		&result.Code,
		&result.Description,
		&result.Details,
		&result.Status,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.ModifiedBy,
		&result.ModifiedAt,
		&result.Vers); err != nil {
		return nil, fmt.Errorf("Failed retrieve product record value, error: %v", err)
	}

	uomRepo := unitofmeasurerepository.NewUnitOfMeasureRepository()
	uoms, err := uomRepo.GetByProduct(ctx, id)
	if err != nil {
		return result, err
	}

	result.UnitOfMeasures = uoms

	fieldRepo := productcustomfieldrepository.NewProductCustomFieldRepository()
	fields, err := fieldRepo.GetByProduct(ctx, id)
	if err != nil {
		return result, err
	}

	result.CustomFields = fields

	return result, nil
}

func (prodRepo *productRepository) GetByCatalogue(ctx context.Context, clgCode string) ([]*productmodel.Product, error) {
	result := make([]*productmodel.Product, 0)

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, clg_code, code, descr, details, status, created_by, created_at, modified_by, modified_at, vers
		FROM products
		WHERE clg_code=$1`)
	if err != nil {
		return result, fmt.Errorf("Failed preparing read product, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, clgCode)
	if err != nil {
		return result, fmt.Errorf("Failed reading product, error: %v", err)
	}
	defer rows.Close()

	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return result, fmt.Errorf("Failed retrieve product record, error: %v", err)
			}
			break
		}

		product := productmodel.NewProduct()
		if err := rows.Scan(
			&product.ID,
			&product.CatalogueCode,
			&product.Code,
			&product.Description,
			&product.Details,
			&product.Status,
			&product.CreatedBy,
			&product.CreatedAt,
			&product.ModifiedBy,
			&product.ModifiedAt,
			&product.Vers); err != nil {
			return result, fmt.Errorf("Failed retrieve product record value, error: %v", err)
		}

		result = append(result, product)
	}

	return result, nil
}

func (prodRepo *productRepository) Create(ctx context.Context, data *productmodel.Product) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO products 
			(clg_code, code, descr, details, status, created_by, created_at, modified_by, modified_at, vers) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing insert product, error: %v", err)
	}
	defer stmt.Close()

	var lastInsertID int64
	err = stmt.QueryRowContext(ctx, data.GetCatalogueCode(), data.GetCode(), data.GetDescription(), data.GetDetails(), data.GetStatus(), data.GetCreatedBy(), data.GetCreatedAt(), data.GetModifiedBy(), data.GetModifiedAt()).Scan(&lastInsertID)
	if err != nil {
		return 0, fmt.Errorf("Failed inserting product, error: %v", err)
	}

	return lastInsertID, nil
}

func (prodRepo *productRepository) Update(ctx context.Context, data *productmodel.Product) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`UPDATE products SET code=$1, descr=$2, details=$3, status=$4, modified_by=$5, modified_at=$6, vers=vers+1 
		WHERE id=$7`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing update product, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetCode(), data.GetDescription(), data.GetDetails(), data.GetStatus(), data.GetModifiedBy(), data.GetModifiedAt(), data.GetID())
	if err != nil {
		return 0, fmt.Errorf("Failed updating product, error: %v", err)
	}

	return result.RowsAffected()
}

func (prodRepo *productRepository) Delete(ctx context.Context, id int64) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM products 
		WHERE id=$1`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing delete product, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("Failed deleting product, error: %v", err)
	}

	return result.RowsAffected()
}

func (prodRepo *productRepository) DeleteByCatalogue(ctx context.Context, clgCode string) error {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM products 
		WHERE clg_code=$1`)
	if err != nil {
		return fmt.Errorf("Failed preparing delete product, error: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, clgCode)
	if err != nil {
		return fmt.Errorf("Failed deleting product, error: %v", err)
	}

	return nil
}
