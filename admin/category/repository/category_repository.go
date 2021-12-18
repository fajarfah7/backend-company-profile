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

func NewCategoryRepository(db *sql.DB) admin.CategoryRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Create(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "INSERT INTO categories(name) values($1)"
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, name)
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

func (rep *Repository) Update(ID int, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "UPDATE categories SET name=$1 WHERE id=$2"
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, name, ID)
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

func (rep *Repository) Delete(ID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "DELETE FROM categories WHERE id=$1"
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, ID)
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
