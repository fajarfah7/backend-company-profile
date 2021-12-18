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

func NewCartRepository(db *sql.DB) user.CartRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Create(userID int64) (*user.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "INSERT INTO carts(user_id, status) VALUES($1, $2)"

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	res, err := stmt.ExecContext(ctx, userID, 0)
	if err != nil {
		return nil, err
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affectedRows == 0 {
		return nil, errors.New("Unknown error occured, affected row is 0")
	}

	qGet := "SELECT id, user_id, status FROM carts WHERE user_id=$1 AND status=$2 ORDER BY id DESC LIMIT 1"

	stmt, err = rep.DB.PrepareContext(ctx, qGet)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var getCart user.Cart
	err = stmt.QueryRowContext(ctx, userID, 0).Scan(&getCart.ID, &getCart.UserID, &getCart.Status)
	if err != nil {
		return nil, err
	}

	return &getCart, nil
}

func (rep *Repository) Update(ID int64, status int, paymentCode string) (*user.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "UPDATE carts SET status=$1, payment_code=$2 WHERE id=$3"

	if status == 3 {
		q = "UPDATE carts SET status=$1, received_at=$2 WHERE id=$3"
	}

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	if status == 2 {
		res, err := stmt.ExecContext(ctx, status, paymentCode, ID)
		if err != nil {
			return nil, err
		}
		affectedRows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if affectedRows == 0 {
			return nil, errors.New("Unknown error occured, affected row is 0")
		}

	} else if status == 3 {
		now := time.Now()
		res, err := stmt.ExecContext(ctx, status, now, ID)
		if err != nil {
			return nil, err
		}
		affectedRows, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}

		if affectedRows == 0 {
			return nil, errors.New("Unknown error occured, affected row is 0")
		}
	}

	qGet := "SELECT id, user_id, status, payment_code, received_at FROM carts WHERE id=$1"

	stmt, err = rep.DB.PrepareContext(ctx, qGet)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var cart user.Cart
	err = stmt.QueryRowContext(ctx, ID).Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.PaymentCode, &cart.ReceivedAt)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (rep *Repository) GetActiveCart(userID int64) (*user.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, user_id, status FROM carts WHERE user_id=$1 AND status=$2 ORDER BY id DESC LIMIT 1"

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var getCart user.Cart

	err = stmt.QueryRowContext(ctx, userID, 0).Scan(&getCart.ID, &getCart.UserID, &getCart.Status)
	if err != nil {
		// error on here indicate that there is no data
		return nil, err
	}

	return &getCart, nil
}

func (rep *Repository) GetByID(ID int64) (*user.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, user_id, status, payment_code FROM carts WHERE id=$1"

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var cart user.Cart
	err = stmt.QueryRowContext(ctx, ID).Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.PaymentCode)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

func (rep *Repository) Carts(userID int64, status, page, limit int, key string) (count int64, cartList []*user.Cart, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	search := ""
	if key != "" {
		search = ` AND payment_code LIKE '%` + key + `%' `
	}

	qCount := fmt.Sprintf(`SELECT COUNT(id) FROM carts WHERE status=$1 AND user_ID=$2 %s`, search)
	stmt, err := rep.DB.PrepareContext(ctx, qCount)
	defer stmt.Close()
	if err != nil {
		return int64(0), nil, err
	}

	type counter struct {
		count int64
	}
	var c counter
	stmt.QueryRowContext(ctx, status, userID).Scan(&c.count)
	if c.count == 0 {
		return int64(0), nil, errors.New("There is no data")
	}

	offset := 0
	if page > 0 {
		offset = (page * limit) - limit
	}

	q := fmt.Sprintf(`SELECT id, user_id, status, payment_code, shipment_code 
	FROM carts 
	WHERE status=$1 AND user_id=$2 %s
	ORDER BY id DESC
	LIMIT $3 OFFSET $4`, search)
	fmt.Println(q)

	stmt, err = rep.DB.PrepareContext(ctx, q)
	if err != nil {
		return int64(0), nil, err
	}

	rows, err := stmt.QueryContext(ctx, status, userID, limit, offset)
	if err != nil {
		return int64(0), nil, err
	}

	var carts []*user.Cart
	for rows.Next() {
		var cart user.Cart
		err := rows.Scan(&cart.ID, &cart.UserID, &cart.Status, &cart.PaymentCode, &cart.ShipmentCode)
		if err != nil {
			return int64(0), nil, err
		}

		carts = append(carts, &cart)
	}

	return int64(c.count), carts, nil
}

func (rep *Repository) Search(userID int64, status int, key string) ([]*user.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `SELECT id, user_id, status, payment_code, shipment_code, received_at FROM carts WHERE user_id=$1 AND status=$2 AND payment_code LIKE '%$3%'`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, userID, status, key)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var carts []*user.Cart
	for rows.Next() {
		var cart user.Cart
		err := rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.Status,
			&cart.PaymentCode,
			&cart.ShipmentCode,
			&cart.ReceivedAt,
		)
		if err != nil {
			return nil, err
		}
		carts = append(carts, &cart)
	}
	return carts, nil
}
