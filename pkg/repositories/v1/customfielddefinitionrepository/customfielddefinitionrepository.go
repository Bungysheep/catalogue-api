package customfielddefinitionrepository

import (
	"context"
	"fmt"
	"time"

	"github.com/bungysheep/catalogue-api/pkg/configs"
	customfielddefinitionmodel "github.com/bungysheep/catalogue-api/pkg/models/v1/customfielddefinition"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
)

// ICustomFieldDefinitionRepository type
type ICustomFieldDefinitionRepository interface {
	GetByID(context.Context, int64) (*customfielddefinitionmodel.CustomFieldDefinition, error)
	GetByCatalogue(context.Context, string) ([]*customfielddefinitionmodel.CustomFieldDefinition, error)
	Create(context.Context, *customfielddefinitionmodel.CustomFieldDefinition) (int64, error)
	Update(context.Context, *customfielddefinitionmodel.CustomFieldDefinition) (int64, error)
	Delete(context.Context, int64) (int64, error)
	DeleteByCatalogue(context.Context, string) error
}

type customFieldDefinitionRepository struct {
}

// NewCustomFieldDefinitionRepository - Create custom field definition repository
func NewCustomFieldDefinitionRepository() ICustomFieldDefinitionRepository {
	return &customFieldDefinitionRepository{}
}

func (fieldDefRepo *customFieldDefinitionRepository) GetByID(ctx context.Context, id int64) (*customfielddefinitionmodel.CustomFieldDefinition, error) {
	result := customfielddefinitionmodel.NewCustomFieldDefinition()

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, clg_code, caption, type, mandatory, created_by, created_at, modified_by, modified_at, vers
		FROM custom_field_definitions 
		WHERE id=?`)
	if err != nil {
		return nil, fmt.Errorf("Failed preparing read custom field definition, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("Failed reading custom field definition, error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("Failed retrieve custom field definition record, error: %v", err)
		}
		return nil, nil
	}

	var createdAt string
	var modifiedAt string
	if err := rows.Scan(
		&result.ID,
		&result.CatalogueCode,
		&result.Caption,
		&result.Type,
		&result.Mandatory,
		&result.CreatedBy,
		&createdAt,
		&result.ModifiedBy,
		&modifiedAt,
		&result.Vers); err != nil {
		return nil, fmt.Errorf("Failed retrieve custom field definition record value, error: %v", err)
	}

	result.CreatedAt, _ = time.Parse(configs.DATEFORMAT, createdAt)
	result.ModifiedAt, _ = time.Parse(configs.DATEFORMAT, modifiedAt)

	return result, nil
}

func (fieldDefRepo *customFieldDefinitionRepository) GetByCatalogue(ctx context.Context, clgCode string) ([]*customfielddefinitionmodel.CustomFieldDefinition, error) {
	result := make([]*customfielddefinitionmodel.CustomFieldDefinition, 0)

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, clg_code, caption, type, mandatory, created_by, created_at, modified_by, modified_at, vers
		FROM custom_field_definitions
		WHERE clg_code=?`)
	if err != nil {
		return result, fmt.Errorf("Failed preparing read custom field definition, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, clgCode)
	if err != nil {
		return result, fmt.Errorf("Failed reading custom field definition, error: %v", err)
	}
	defer rows.Close()

	var createdAt string
	var modifiedAt string
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return result, fmt.Errorf("Failed retrieve custom field definition record, error: %v", err)
			}
			break
		}

		fieldDef := customfielddefinitionmodel.NewCustomFieldDefinition()
		if err := rows.Scan(
			&fieldDef.ID,
			&fieldDef.CatalogueCode,
			&fieldDef.Caption,
			&fieldDef.Type,
			&fieldDef.Mandatory,
			&fieldDef.CreatedBy,
			&createdAt,
			&fieldDef.ModifiedBy,
			&modifiedAt,
			&fieldDef.Vers); err != nil {
			return result, fmt.Errorf("Failed retrieve custom field definition record value, error: %v", err)
		}

		fieldDef.CreatedAt, _ = time.Parse(configs.DATEFORMAT, createdAt)
		fieldDef.ModifiedAt, _ = time.Parse(configs.DATEFORMAT, modifiedAt)

		result = append(result, fieldDef)
	}

	return result, nil
}

func (fieldDefRepo *customFieldDefinitionRepository) Create(ctx context.Context, data *customfielddefinitionmodel.CustomFieldDefinition) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO custom_field_definitions 
			(clg_code, caption, type, mandatory, created_by, created_at, modified_by, modified_at, vers) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, 1)`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing insert custom field definition, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetCatalogueCode(), data.GetCaption(), data.GetType(), data.GetMandatory(), data.GetCreatedBy(), data.GetCreatedAt(), data.GetModifiedBy(), data.GetModifiedAt())
	if err != nil {
		return 0, fmt.Errorf("Failed inserting custom field definition, error: %v", err)
	}

	return result.RowsAffected()
}

func (fieldDefRepo *customFieldDefinitionRepository) Update(ctx context.Context, data *customfielddefinitionmodel.CustomFieldDefinition) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`UPDATE custom_field_definitions SET caption=?, type=?, mandatory=?, modified_by=?, modified_at=?, vers=vers+1 
		WHERE id=?`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing update custom field definition, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetCaption(), data.GetType(), data.GetMandatory(), data.GetModifiedBy(), data.GetModifiedAt(), data.GetID())
	if err != nil {
		return 0, fmt.Errorf("Failed updating custom field definition, error: %v", err)
	}

	return result.RowsAffected()
}

func (fieldDefRepo *customFieldDefinitionRepository) Delete(ctx context.Context, id int64) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM custom_field_definitions 
		WHERE id=?`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing delete custom field definition, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("Failed deleting custom field definition, error: %v", err)
	}

	return result.RowsAffected()
}

func (fieldDefRepo *customFieldDefinitionRepository) DeleteByCatalogue(ctx context.Context, clgCode string) error {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`DELETE FROM custom_field_definitions 
		WHERE clg_code=?`)
	if err != nil {
		return fmt.Errorf("Failed preparing delete v, error: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, clgCode)
	if err != nil {
		return fmt.Errorf("Failed deleting custom field definition, error: %v", err)
	}

	return nil
}
