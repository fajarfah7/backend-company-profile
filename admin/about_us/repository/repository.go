package repository

import (
	"context"
	"database/sql"
	"errors"
	"project_company_profile/models/admin"
	"time"
)

// Repository struct which will have methods CreateOrUpdate and GetAll
type Repository struct {
	DB *sql.DB
}

// NewAboutUsRepository return admin.NewAboutUsRepository interface
func NewAboutUsRepository(db *sql.DB) admin.AboutUsRepository {
	return &Repository{
		DB: db,
	}
}

// CreateOrUpdate will be used to store or update, if there is no data will create new data, if there was data will do update
func (rep *Repository) CreateOrUpdate(textName, textType, newText string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `select id, coalesce(image_id, 0), name, type, text, created_at, updated_at from texts where type = $1`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	row := stmt.QueryRowContext(ctx, textType)

	var aboutUs admin.AboutUs
	_ = row.Scan(
		&aboutUs.ID,
		&aboutUs.ImageID,
		&aboutUs.Name,
		&aboutUs.Type,
		&aboutUs.Text,
		&aboutUs.CreatedAt,
		&aboutUs.UpdatedAt,
	)
	// if err != nil {
	// 	return err
	// }

	if aboutUs.ID == 0 {
		q := `insert into texts (name, type, text) values($1, $2, $3)`

		stmt, err := rep.DB.PrepareContext(ctx, q)
		defer stmt.Close()
		if err != nil {
			return err
		}

		res, err := stmt.ExecContext(ctx, textName, textType, newText)
		if err != nil {
			return err
		}

		affectedRows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("Failed set data")
		}

	} else {
		q := `update texts set text = $1 where id = $2`

		stmt, err := rep.DB.PrepareContext(ctx, q)
		defer stmt.Close()
		if err != nil {
			return err
		}

		res, err := stmt.ExecContext(ctx, newText, aboutUs.ID)
		if err != nil {
			return err
		}

		affectedRows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affectedRows == 0 {
			return errors.New("Failed set data")
		}
	}

	return nil
}

// GetAll return all about us data from database which has type like 'about_us_%'
func (rep *Repository) GetAll() ([]*admin.AboutUs, error) {
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

	var listAboutUs []*admin.AboutUs

	for rows.Next() {
		aboutUs := new(admin.AboutUs)

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
