package cataloguerepository

import (
	"context"
	"fmt"

	cataloguemodel "github.com/bungysheep/catalogue-api/pkg/models/v1/catalogue"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
	"github.com/bungysheep/catalogue-api/pkg/repositories/v1/customfielddefinitionrepository"
)

// ICatalogueRepository type
type ICatalogueRepository interface {
	GetByID(context.Context, string) (*cataloguemodel.Catalogue, error)
	GetAll(context.Context) ([]*cataloguemodel.Catalogue, error)
	Create(context.Context, *cataloguemodel.Catalogue) (int64, error)
	Update(context.Context, *cataloguemodel.Catalogue) (int64, error)
	Delete(context.Context, string) (int64, error)
}

type catalogueRepository struct {
}

// NewCatalogueRepository - Create catalogue repository
func NewCatalogueRepository() ICatalogueRepository {
	return &catalogueRepository{}
}

func (clgRepo *catalogueRepository) GetByID(ctx context.Context, code string) (*cataloguemodel.Catalogue, error) {
	result := cataloguemodel.NewCatalogue()

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT code, descr, details, status, created_by, created_at, modified_by, modified_at, vers
		FROM catalogues 
		WHERE code=$1`)
	if err != nil {
		return nil, fmt.Errorf("Failed preparing read catalogue, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("Failed reading catalogue, error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("Failed retrieve catalogue record, error: %v", err)
		}
		return nil, nil
	}

	if err := rows.Scan(
		&result.Code,
		&result.Description,
		&result.Details,
		&result.Status,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.ModifiedBy,
		&result.ModifiedAt,
		&result.Vers); err != nil {
		return nil, fmt.Errorf("Failed retrieve catalogue record value, error: %v", err)
	}

	fieldDefRepo := customfielddefinitionrepository.NewCustomFieldDefinitionRepository()
	fieldDefs, err := fieldDefRepo.GetByCatalogue(ctx, code)
	if err != nil {
		return result, err
	}

	result.CustomFieldDefinitions = fieldDefs

	return result, nil
}

func (clgRepo *catalogueRepository) GetAll(ctx context.Context) ([]*cataloguemodel.Catalogue, error) {
	result := make([]*cataloguemodel.Catalogue, 0)

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT code, descr, details, status, created_by, created_at, modified_by, modified_at, vers
		FROM catalogues`)
	if err != nil {
		return result, fmt.Errorf("Failed preparing read catalogue, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed reading catalogue, error: %v", err)
	}
	defer rows.Close()

	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return result, fmt.Errorf("Failed retrieve catalogue record, error: %v", err)
			}
			break
		}

		catalogue := cataloguemodel.NewCatalogue()
		if err := rows.Scan(
			&catalogue.Code,
			&catalogue.Description,
			&catalogue.Details,
			&catalogue.Status,
			&catalogue.CreatedBy,
			&catalogue.CreatedAt,
			&catalogue.ModifiedBy,
			&catalogue.ModifiedAt,
			&catalogue.Vers); err != nil {
			return result, fmt.Errorf("Failed retrieve catalogue record value, error: %v", err)
		}

		result = append(result, catalogue)
	}

	return result, nil
}

func (clgRepo *catalogueRepository) Create(ctx context.Context, data *cataloguemodel.Catalogue) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO catalogues 
			(code, descr, details, status, created_by, created_at, modified_by, modified_at, vers) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 1)`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing insert catalogue, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetCode(), data.GetDescription(), data.GetDetails(), data.GetStatus(), data.GetCreatedBy(), data.GetCreatedAt(), data.GetModifiedBy(), data.GetModifiedAt())
	if err != nil {
		return 0, fmt.Errorf("Failed inserting catalogue, error: %v", err)
	}

	return result.RowsAffected()
}

func (clgRepo *catalogueRepository) Update(ctx context.Context, data *cataloguemodel.Catalogue) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`UPDATE catalogues SET descr=$1, details=$2, status=$3, modified_by=$4, modified_at=$5, vers=vers+1 
		WHERE code=$6`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing update catalogue, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetDescription(), data.GetDetails(), data.GetStatus(), data.GetModifiedBy(), data.GetModifiedAt(), data.GetCode())
	if err != nil {
		return 0, fmt.Errorf("Failed updating catalogue, error: %v", err)
	}

	return result.RowsAffected()
}

func (clgRepo *catalogueRepository) Delete(ctx context.Context, code string) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM catalogues 
		WHERE code=$1`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing delete catalogue, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, code)
	if err != nil {
		return 0, fmt.Errorf("Failed deleting catalogue, error: %v", err)
	}

	return result.RowsAffected()
}
