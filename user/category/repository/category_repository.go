package repository

import (
	"context"
	"database/sql"
	"project_company_profile/models/user"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewCategoryRepository(db *sql.DB) user.CategoryRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Category(ID int) (*user.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, name FROM categories WHERE id=$1"
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, ID)
	var category user.Category
	row.Scan(&category.ID, &category.Name)

	return &category, nil
}

func (rep *Repository) Categories() ([]*user.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, name FROM categories ORDER BY id DESC"
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var categories []*user.Category
	for rows.Next() {
		var category user.Category
		err := rows.Scan(&category.ID, &category.Name)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	return categories, nil
}
