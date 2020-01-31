package unitofmeasurerepository

import (
	"context"
	"fmt"

	unitofmeasuremodel "github.com/bungysheep/catalogue-api/pkg/models/v1/unitofmeasure"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
)

// IUnitOfMeasureRepository type
type IUnitOfMeasureRepository interface {
	GetByID(context.Context, int64) (*unitofmeasuremodel.UnitOfMeasure, error)
	GetByProduct(context.Context, int64) ([]*unitofmeasuremodel.UnitOfMeasure, error)
	Create(context.Context, *unitofmeasuremodel.UnitOfMeasure) (int64, error)
	Update(context.Context, *unitofmeasuremodel.UnitOfMeasure) (int64, error)
	Delete(context.Context, int64) (int64, error)
}

type unitOfMeasureRepository struct {
}

// NewUnitOfMeasureRepository - Create unit of measure repository
func NewUnitOfMeasureRepository() IUnitOfMeasureRepository {
	return &unitOfMeasureRepository{}
}

func (uomRepo *unitOfMeasureRepository) GetByID(ctx context.Context, id int64) (*unitofmeasuremodel.UnitOfMeasure, error) {
	result := unitofmeasuremodel.NewUnitOfMeasure()

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, prod_id, code, descr, is_default, ratio, vers
		FROM product_uoms 
		WHERE id=?`)
	if err != nil {
		return nil, fmt.Errorf("Failed preparing read unit of measure, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Failed reading unit of measure, error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("Failed retrieve unit of measure record, error: %v", err)
		}
		return nil, nil
	}

	if err := rows.Scan(
		&result.ID,
		&result.ProdID,
		&result.Code,
		&result.Description,
		&result.IsDefault,
		&result.Ratio,
		&result.Vers); err != nil {
		return nil, fmt.Errorf("Failed retrieve unit of measure record value, error: %v", err)
	}

	return result, nil
}

func (uomRepo *unitOfMeasureRepository) GetByProduct(ctx context.Context, prodID int64) ([]*unitofmeasuremodel.UnitOfMeasure, error) {
	result := make([]*unitofmeasuremodel.UnitOfMeasure, 0)

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, prod_id, code, descr, is_default, ratio, vers
		FROM product_uoms
		WHERE prod_id=?
		ORDER BY is_default DESC, ratio ASC`)
	if err != nil {
		return result, fmt.Errorf("Failed preparing read unit of measure, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, prodID)
	if err != nil {
		return result, fmt.Errorf("Failed reading unit of measure, error: %v", err)
	}
	defer rows.Close()

	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return result, fmt.Errorf("Failed retrieve unit of measure record, error: %v", err)
			}
			break
		}

		uom := unitofmeasuremodel.NewUnitOfMeasure()
		if err := rows.Scan(
			&uom.ID,
			&uom.ProdID,
			&uom.Code,
			&uom.Description,
			&uom.IsDefault,
			&uom.Ratio,
			&uom.Vers); err != nil {
			return result, fmt.Errorf("Failed retrieve unit of measure record value, error: %v", err)
		}

		result = append(result, uom)
	}

	return result, nil
}

func (uomRepo *unitOfMeasureRepository) Create(ctx context.Context, data *unitofmeasuremodel.UnitOfMeasure) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO catalogues 
			(prod_id, code, descr, is_default, ratio, vers) 
		VALUES (?, ?, ?, ?, ?, 1)`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing insert unit or measure, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetProdID(), data.GetCode(), data.GetDescription(), data.GetIsDefault(), data.GetRatio())
	if err != nil {
		return 0, fmt.Errorf("Failed inserting unit or measure, error: %v", err)
	}

	return result.RowsAffected()
}

func (uomRepo *unitOfMeasureRepository) Update(ctx context.Context, data *unitofmeasuremodel.UnitOfMeasure) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`UPDATE product_uoms SET code=?, descr=?, is_default=?, ratio=?, vers=vers+1 
		WHERE id=?`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing update unit or measure, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetCode(), data.GetDescription(), data.GetIsDefault(), data.GetRatio(), data.GetID())
	if err != nil {
		return 0, fmt.Errorf("Failed updating unit or measure, error: %v", err)
	}

	return result.RowsAffected()
}

func (uomRepo *unitOfMeasureRepository) Delete(ctx context.Context, id int64) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM product_uoms 
		WHERE id=?`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing delete unit or measure, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("Failed deleting unit or measure, error: %v", err)
	}

	return result.RowsAffected()
}
