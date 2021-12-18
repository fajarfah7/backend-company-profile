package repository

import (
	"context"
	"database/sql"
	"fmt"
	"project_company_profile/models/user"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewTextRepository(db *sql.DB) user.TextRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Texts(textType string) ([]*user.Text, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := fmt.Sprintf(`SELECT id, coalesce(image_id, 0), name, type, text, created_at, updated_at FROM texts WHERE type LIKE '%%%s%%'`, textType)
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

	var texts []*user.Text
	for rows.Next() {
		text := new(user.Text)

		err := rows.Scan(
			&text.ID,
			&text.ImageID,
			&text.Name,
			&text.Type,
			&text.Text,
			&text.CreatedAt,
			&text.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		texts = append(texts, text)
	}

	return texts, nil
}
