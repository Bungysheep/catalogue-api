package productcustomfieldrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/configs"
	productcustomfieldmodel "github.com/bungysheep/catalogue-api/pkg/models/v1/productcustomfield"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
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
	result := productcustomfieldmodel.NewProductCustomField()

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, prod_id, field_id, alpha_value, numeric_value, date_value
		FROM product_custom_fields 
		WHERE id=?`)
	if err != nil {
		return nil, fmt.Errorf("Failed preparing read product custom field, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Failed reading product custom field, error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("Failed retrieve product custom field record, error: %v", err)
		}
		return nil, nil
	}

	var dateValue string
	if err := rows.Scan(
		&result.ID,
		&result.ProdID,
		&result.FieldID,
		&result.AlphaValue,
		&result.NumericValue,
		&dateValue); err != nil {
		return nil, fmt.Errorf("Failed retrieve product custom field record value, error: %v", err)
	}

	result.DateValue, _ = time.Parse(configs.DATEFORMAT, dateValue)

	return result, nil
}

func (pcfRepo *productCustomFieldRepository) GetByProduct(ctx context.Context, prodID int64) ([]*productcustomfieldmodel.ProductCustomField, error) {
	result := make([]*productcustomfieldmodel.ProductCustomField, 0)

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, prod_id, field_id, alpha_value, numeric_value, date_value
		FROM product_custom_fields
		WHERE prod_id=?
		ORDER BY field_id, id ASC`)
	if err != nil {
		return result, fmt.Errorf("Failed preparing read product custom field, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, prodID)
	if err != nil {
		return result, fmt.Errorf("Failed reading product custom field, error: %v", err)
	}
	defer rows.Close()

	var dateValue string
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return result, fmt.Errorf("Failed retrieve product custom field record, error: %v", err)
			}
			break
		}

		field := productcustomfieldmodel.NewProductCustomField()
		if err := rows.Scan(
			&field.ID,
			&field.ProdID,
			&field.FieldID,
			&field.AlphaValue,
			&field.NumericValue,
			&dateValue); err != nil {
			return result, fmt.Errorf("Failed retrieve product custom field record value, error: %v", err)
		}

		field.DateValue, _ = time.Parse(configs.DATEFORMAT, dateValue)

		result = append(result, field)
	}

	return result, nil
}

func (pcfRepo *productCustomFieldRepository) Create(ctx context.Context, data *productcustomfieldmodel.ProductCustomField) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO product_custom_fields 
			(prod_id, field_id, alpha_value, numeric_value, date_value) 
		VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing insert product custom field, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetProdID(), data.GetFieldID(), data.GetAlphaValue(), data.GetNumericValue(), data.GetDateValue())
	if err != nil {
		return 0, fmt.Errorf("Failed inserting product custom field, error: %v", err)
	}

	return result.RowsAffected()
}

func (pcfRepo *productCustomFieldRepository) Update(ctx context.Context, data *productcustomfieldmodel.ProductCustomField) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`UPDATE product_custom_fields SET alpha_value=?, numeric_value=?, date_value=? 
		WHERE id=?`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing update product custom field, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetAlphaValue(), data.GetNumericValue(), data.GetDateValue(), data.GetID(), data.GetID())
	if err != nil {
		return 0, fmt.Errorf("Failed updating product custom field, error: %v", err)
	}

	return result.RowsAffected()
}

func (pcfRepo *productCustomFieldRepository) Delete(ctx context.Context, id int64) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM product_custom_fields 
		WHERE id=?`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing delete product custom field, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("Failed deleting product custom field, error: %v", err)
	}

	return result.RowsAffected()
}

func (pcfRepo *productCustomFieldRepository) DeleteByProduct(ctx context.Context, prodID int64) error {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM product_custom_fields 
		WHERE prod_id=?`)
	if err != nil {
		return fmt.Errorf("Failed preparing delete v, error: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, prodID)
	if err != nil {
		return fmt.Errorf("Failed deleting product custom field, error: %v", err)
	}

	return nil
}
