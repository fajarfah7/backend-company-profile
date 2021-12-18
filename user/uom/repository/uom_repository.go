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

func NewUomRepository(db *sql.DB) user.UomRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Uom(ID int) (*user.Uom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, name, amount FROM uoms WHERE id=$1"
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	row := stmt.QueryRowContext(ctx, ID)
	var uom user.Uom
	row.Scan(&uom.ID, &uom.Name, &uom.Amount)

	return &uom, nil
}

func (rep *Repository) Uoms() ([]*user.Uom, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, name, amount FROM uoms"
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

	var uoms []*user.Uom
	for rows.Next() {
		var uom user.Uom
		err := rows.Scan(&uom.ID, &uom.Name, &uom.Amount)
		if err != nil {
			return nil, err
		}
		uoms = append(uoms, &uom)
	}
	return uoms, nil
}
