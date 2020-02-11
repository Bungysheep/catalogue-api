package userrepository

import (
	"context"
	"fmt"

	usermodel "github.com/bungysheep/catalogue-api/pkg/models/v1/user"
	"github.com/bungysheep/catalogue-api/pkg/protocols/database"
)

// IUserRepository type
type IUserRepository interface {
	GetAll(context.Context) ([]*usermodel.User, error)
	GetByUsername(context.Context, string) (*usermodel.User, error)
	Create(context.Context, *usermodel.User) (int64, error)
}

type userRepository struct {
}

// NewUserRepository - Create user repository
func NewUserRepository() IUserRepository {
	return &userRepository{}
}

func (usrRepo *userRepository) GetAll(ctx context.Context) ([]*usermodel.User, error) {
	result := make([]*usermodel.User, 0)

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT username, name, email, password, status, created_by, created_at, modified_by, modified_at, vers
		FROM users`)
	if err != nil {
		return result, fmt.Errorf("Failed preparing read user, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return result, fmt.Errorf("Failed reading user, error: %v", err)
	}
	defer rows.Close()

	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return result, fmt.Errorf("Failed retrieve catalogue record, error: %v", err)
			}
			break
		}

		user := usermodel.NewUser()
		if err := rows.Scan(
			&user.Username,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Status,
			&user.CreatedBy,
			&user.CreatedAt,
			&user.ModifiedBy,
			&user.ModifiedAt,
			&user.Vers); err != nil {
			return result, fmt.Errorf("Failed retrieve user record value, error: %v", err)
		}

		result = append(result, user)
	}

	return result, nil
}

func (usrRepo *userRepository) GetByUsername(ctx context.Context, username string) (*usermodel.User, error) {
	result := usermodel.NewUser()

	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`SELECT username, name, email, password, status, created_by, created_at, modified_by, modified_at, vers
		FROM users
		WHERE username=$1`)
	if err != nil {
		return nil, fmt.Errorf("Failed preparing read user, error: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("Failed reading user, error: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("Failed retrieve user record, error: %v", err)
		}
		return nil, nil
	}

	if err := rows.Scan(
		&result.Username,
		&result.Name,
		&result.Email,
		&result.Password,
		&result.Status,
		&result.CreatedBy,
		&result.CreatedAt,
		&result.ModifiedBy,
		&result.ModifiedAt,
		&result.Vers); err != nil {
		return nil, fmt.Errorf("Failed retrieve user record value, error: %v", err)
	}

	return result, nil
}

func (usrRepo *userRepository) Create(ctx context.Context, data *usermodel.User) (int64, error) {
	conn, err := database.DbConnection.Conn(ctx)
	if err != nil {
		return 0, fmt.Errorf("Failed connecting to database, error: %v", err)
	}
	defer conn.Close()

	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO users (username, name, email, password, status, created_by, created_at, modified_by, modified_at, vers) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, 1)`)
	if err != nil {
		return 0, fmt.Errorf("Failed preparing create user, error: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, data.GetUsername(), data.GetName(), data.GetEmail(), data.GetPassword(), data.GetStatus(), data.GetCreatedBy(), data.GetCreatedAt(), data.GetModifiedBy(), data.GetModifiedAt())
	if err != nil {
		return 0, fmt.Errorf("Failed inserting user, error: %v", err)
	}

	return result.RowsAffected()
}
