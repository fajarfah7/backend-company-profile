package repository

import (
	"context"
	"database/sql"
	"errors"
	"project_company_profile/models/admin"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewAdminUserRepository(db *sql.DB) admin.UserRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) GetByID(ID int64) (*admin.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `SELECT 
		id,
		name,
		email,
		username,
		password,
		address,
		phone_number,
		created_at,
		updated_at
	FROM users
	WHERE id = $1`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var user admin.User
	row := stmt.QueryRowContext(ctx, ID)
	row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Address,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if user.ID == 0 {
		return nil, errors.New("Unknown error occured, failed get user data")
	}

	return &user, nil
}
