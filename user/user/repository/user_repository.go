package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"project_company_profile/models/user"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) user.UserRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Create(user user.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `INSERT INTO users(
		name,
		email,
		username,
		password,
		address,
		phone_number,
		created_at,
		updated_at
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx,
		user.Name,
		user.Email,
		user.Username,
		user.Password,
		user.Address,
		user.PhoneNumber,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return errors.New("Unknown error occured, affected row is 0")
	}

	return nil

}

func (rep *Repository) Update(user user.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `UPDATE users SET
		name = $1,
		email = $2,
		address = $3,
		phone_number = $4,
		updated_at = $5
	WHERE id = $6`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx,
		user.Name,
		user.Email,
		user.Address,
		user.PhoneNumber,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return err
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return errors.New("Unknown error occured, affected row is 0")
	}

	return nil
}

func (rep *Repository) CheckPassword(userID int64) (hashedPassword *string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `SELECT password FROM users WHERE id = $1`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRowContext(ctx, userID)

	row.Scan(&hashedPassword)

	fmt.Println(hashedPassword)

	return hashedPassword, nil
}

func (rep *Repository) UpdatePassword(userID int64, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "UPDATE users SET password=$1 WHERE id=$2"

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, newPassword, userID)
	if err != nil {
		return err
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return errors.New("Unknown error occured, affected row is 0")
	}

	return nil
}

func (rep *Repository) Delete(userID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// stmtCarts := `DELETE FROM carts WHERE user_id = $1`
	// _, err = rep.DB.ExecContext(ctx, stmtCarts, userID)
	// if err != nil {
	// 	defer rep.DB.Close()
	// 	return err
	// }

	q := `DELETE FROM users WHERE id = $1`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx, userID)
	if err != nil {
		return err
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return errors.New("Unknown error occured, affected row is 0")
	}

	return nil
}

func (rep *Repository) GetByUsername(username string) (*user.User, error) {
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
	WHERE username = $1`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var user user.User
	row := stmt.QueryRowContext(ctx, username)

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
		return nil, errors.New("Wrong username or password")
	}

	return &user, nil
}
