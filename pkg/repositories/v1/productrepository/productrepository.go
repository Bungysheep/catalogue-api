package productrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/configs"
	productmodel "github.com/bungysheep/catalogue-api/pkg/models/v1/product"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
)

// IProductRepository type
type IProductRepository interface {
	GetByID(context.Context, int64) (*productmodel.Product, error)
	GetByCatalogue(context.Context, string) ([]*productmodel.Product, error)
	Create(context.Context, *productmodel.Product) (int64, error)
	Update(context.Context, *productmodel.Product) (int64, error)
	Delete(context.Context, int64) (int64, error)
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
		WHERE id=?`)
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

	var createdAt string
	var modifiedAt string
	if err := rows.Scan(
		&result.ID,
		&result.CatalogueCode,
		&result.Code,
		&result.Description,
		&result.Details,
		&result.Status,
		&result.CreatedBy,
		&createdAt,
		&result.ModifiedBy,
		&modifiedAt,
		&result.Vers); err != nil {
		return nil, fmt.Errorf("Failed retrieve product record value, error: %v", err)
	}

	result.CreatedAt, _ = time.Parse(configs.DATEFORMAT, createdAt)
	result.ModifiedAt, _ = time.Parse(configs.DATEFORMAT, modifiedAt)

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
		WHERE clg_code=?`)
	if err != nil {
		return result, fmt.Errorf("Failed preparing read product, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, clgCode)
	if err != nil {
		return result, fmt.Errorf("Failed reading product, error: %v", err)
	}
	defer rows.Close()

	var createdAt string
	var modifiedAt string
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
			&createdAt,
			&product.ModifiedBy,
			&modifiedAt,
			&product.Vers); err != nil {
			return result, fmt.Errorf("Failed retrieve product record value, error: %v", err)
		}

		product.CreatedAt, _ = time.Parse(configs.DATEFORMAT, createdAt)
		product.ModifiedAt, _ = time.Parse(configs.DATEFORMAT, modifiedAt)

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
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 1)`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing insert product, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetCatalogueCode(), data.GetCode(), data.GetDescription(), data.GetDetails(), data.GetStatus(), data.GetCreatedBy(), data.GetCreatedAt(), data.GetModifiedBy(), data.GetModifiedAt())
	if err != nil {
		return 0, fmt.Errorf("Failed inserting product, error: %v", err)
	}

	return result.LastInsertId()
}

func (prodRepo *productRepository) Update(ctx context.Context, data *productmodel.Product) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`UPDATE products SET code=?, descr=?, details=?, status=?, modified_by=?, modified_at=?, vers=vers+1 
		WHERE id=?`)
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
		WHERE id=?`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing delete v, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("Failed deleting product, error: %v", err)
	}

	return result.RowsAffected()
}
