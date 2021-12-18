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

func NewAboutUsRepository(db *sql.DB) user.AboutUsRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) GetAll() ([]*user.AboutUs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `select id, coalesce(image_id, 0), name, type, text, created_at, updated_at from texts where type like 'about_us_%'`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx)
	defer rows.Close()
	if err != nil {
		rep.DB.Close()
		return nil, err
	}

	var listAboutUs []*user.AboutUs

	for rows.Next() {
		aboutUs := new(user.AboutUs)

		err := rows.Scan(
			&aboutUs.ID,
			&aboutUs.ImageID,
			&aboutUs.Name,
			&aboutUs.Type,
			&aboutUs.Text,
			&aboutUs.CreatedAt,
			&aboutUs.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		listAboutUs = append(listAboutUs, aboutUs)
	}

	return listAboutUs, nil
}
