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

func NewCartRepository(db *sql.DB) admin.CartRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Cart(ID int64) (*admin.Cart, error) {
	return nil, errors.New("Not ready yet")
}

func (rep *Repository) Carts(status, page, limit int) (totalData int64, data []*admin.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	qCount := `SELECT COUNT(id) FROM carts WHERE status=$1`
	stmt, err := rep.DB.PrepareContext(ctx, qCount)
	defer stmt.Close()
	if err != nil {
		return int64(0), nil, err
	}

	type counter struct {
		count int64
	}
	var c counter
	stmt.QueryRowContext(ctx, status).Scan(&c.count)
	if c.count == 0 {
		return int64(0), nil, errors.New("There is no data")
	}

	offset := 0
	if page > 0 {
		offset = (page * limit) - limit
	}

	q := `SELECT id, user_id, status, payment_code, shipment_code 
		FROM carts 
		WHERE status=$1 
		ORDER BY id DESC
		LIMIT $2 OFFSET $3`
	stmt, err = rep.DB.PrepareContext(ctx, q)
	if err != nil {
		return int64(0), nil, err
	}

	rows, err := stmt.QueryContext(ctx, status, limit, offset)
	if err != nil {
		return int64(0), nil, err
	}

	var carts []*admin.Cart
	for rows.Next() {
		var cart admin.Cart
		err := rows.Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.PaymentCode, &cart.ShipmentCode)
		if err != nil {
			return int64(0), nil, err
		}

		carts = append(carts, &cart)
	}

	return int64(c.count), carts, nil
}

func (rep *Repository) Update(ID int64, status int, shipmentCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `UPDATE carts SET status=$1, shipment_code=$2 WHERE id=$3`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, status, shipmentCode, ID)
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
