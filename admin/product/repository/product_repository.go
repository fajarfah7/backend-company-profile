package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"project_company_profile/database"
	"project_company_profile/models/admin"
	"time"
)

// Repository struct which will have methods Create, Update, Delete, GetByID and GetAll
type Repository struct {
	DB *sql.DB
}

// NewProductRepository return admin.NewProductRepository interface
func NewProductRepository(db *sql.DB) admin.ProductRepository {
	return &Repository{
		DB: db,
	}
}

// Create is used to store new data to database
func (rep *Repository) Create(prod admin.Product) error {
	db, err := database.OpenDB()
	defer db.Close()
	if err != nil {
		return err
	}
	rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `INSERT INTO products (
		uom_id,
		name,
		description,
		dimension_description,
		stock,
		price,
		discount,
		final_price,
		is_have_expiry,
		expired_at,
		product_image,
		status,
		category_id,
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
		$8,
		$9,
		$10,
		$11,
		$12,
		$13,
		$14,
		$15
	)`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx,
		prod.UomID,
		prod.Name,
		prod.Description,
		prod.DimensionDescription,
		prod.Stock,
		prod.Price,
		prod.Discount,
		prod.FinalPrice,
		prod.IsHaveExpiry,
		prod.ExpiredAt,
		prod.ProductImage,
		prod.Status,
		prod.CategoryID,
		prod.CreatedAt,
		prod.UpdatedAt,
	)
	if err != nil {
		return err
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow == 0 {
		return errors.New("Unkown error, affected row is 0")
	}

	return nil
}

// Update is used to update product data on database
func (rep *Repository) Update(prod admin.Product) error {
	db, err := database.OpenDB()
	defer db.Close()
	if err != nil {
		return err
	}

	rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `UPDATE products SET 
		uom_id 					= $1,
		name 					= $2,
		description 			= $3,
		dimension_description 	= $4,
		stock 					= $5,
		price 					= $6,
		discount 				= $7,
		final_price 			= $8,
		is_have_expiry 			= $9,
		expired_at 				= $10,
		product_image 			= $11,
		status 					= $12,
		category_id				= $13,
		updated_at 				= $14
	WHERE id = $15`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}
	res, err := stmt.ExecContext(ctx,
		prod.UomID,
		prod.Name,
		prod.Description,
		prod.DimensionDescription,
		prod.Stock,
		prod.Price,
		prod.Discount,
		prod.FinalPrice,
		prod.IsHaveExpiry,
		prod.ExpiredAt,
		prod.ProductImage,
		prod.Status,
		prod.CategoryID,
		prod.UpdatedAt,
		prod.ID,
	)
	if err != nil {
		defer rep.DB.Close()
		return err
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		defer rep.DB.Close()
		return err
	}

	if affectedRow == 0 {
		defer rep.DB.Close()
		return errors.New("Unkown error, affected row is 0")
	}

	defer rep.DB.Close()
	return nil
}

// Delete will delete product data on database
func (rep *Repository) Delete(id int64) error {
	db, err := database.OpenDB()
	defer db.Close()
	if err != nil {
		return err
	}

	rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `DELETE FROM products WHERE id = $1`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	affectedRow, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRow != 1 {
		msg := fmt.Sprintf("affected data is %d", affectedRow)
		return errors.New(msg)
	}

	return nil
}

// GetByID find the data based on given ID
func (rep *Repository) GetByID(id int64) (*admin.Product, error) {
	// open db connection
	db, err := database.OpenDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `SELECT 
		id,
		uom_id,
		category_id,
		name,
		description,
		dimension_description,
		stock,
		price,
		discount,
		final_price,
		is_have_expiry,
		expired_at,
		product_image,
		status,
		created_at,
		updated_at 
	FROM products 
	WHERE id = $1`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var prod admin.Product
	row := stmt.QueryRowContext(ctx, id)
	row.Scan(
		&prod.ID,
		&prod.UomID,
		&prod.CategoryID,
		&prod.Name,
		&prod.Description,
		&prod.DimensionDescription,
		&prod.Stock,
		&prod.Price,
		&prod.Discount,
		&prod.FinalPrice,
		&prod.IsHaveExpiry,
		&prod.ExpiredAt,
		&prod.ProductImage,
		&prod.Status,
		&prod.CreatedAt,
		&prod.UpdatedAt,
	)
	if prod.ID == 0 {
		return nil, errors.New("data not found")
	}

	return &prod, nil
}

// GetProductByID find the data based on given ID, will return *admin.Product and nil
func (rep *Repository) GetProductByID(id int64) (*admin.Product, error) {
	// open db connection
	db, err := database.OpenDB()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `SELECT 
		id,
		uom_id,
		category_id,
		name,
		description,
		dimension_description,
		stock,
		price,
		discount,
		final_price,
		is_have_expiry,
		expired_at,
		product_image,
		status,
		created_at,
		updated_at 
	FROM products 
	WHERE id = $1`

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var prod admin.Product
	row := stmt.QueryRowContext(ctx, id)
	row.Scan(
		&prod.ID,
		&prod.UomID,
		&prod.CategoryID,
		&prod.Name,
		&prod.Description,
		&prod.DimensionDescription,
		&prod.Stock,
		&prod.Price,
		&prod.Discount,
		&prod.FinalPrice,
		&prod.IsHaveExpiry,
		&prod.ExpiredAt,
		&prod.ProductImage,
		&prod.Status,
		&prod.CreatedAt,
		&prod.UpdatedAt,
	)
	if prod.ID == 0 {
		return nil, errors.New("data not found")
	}

	return &prod, nil
}

// GetAll return data based on as much as the limit(used for pagination purpose)
func (rep *Repository) GetAll(page, limit int) (int64, []*admin.Product, error) {
	// open db connection
	db, err := database.OpenDB()
	defer db.Close()
	if err != nil {
		return 0, nil, err
	}

	rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// set limit and offset for pagination purpose
	lim := fmt.Sprintf(" LIMIT %d ", limit)
	off := ""
	if page > 1 {
		offsetCounter := (page * limit) - limit
		off = fmt.Sprintf(" OFFSET %d ", offsetCounter)
	}

	// statement to count all data
	type localCounter struct {
		count int64
	}
	var counter localCounter
	qCount := `SELECT COUNT(id) FROM products`
	stmt, err := rep.DB.PrepareContext(ctx, qCount)
	defer stmt.Close()
	if err != nil {
		return 0, nil, err
	}

	count := stmt.QueryRowContext(ctx)
	count.Scan(&counter.count)
	if counter.count == 0 {
		return 0, nil, errors.New("There is no data")
	}

	// statement to do query
	var prods []*admin.Product
	q := fmt.Sprintf(`SELECT 
		id,
		uom_id,
		category_id,
		name,
		description,
		dimension_description,
		stock,
		price,
		discount,
		final_price,
		is_have_expiry,
		expired_at,
		product_image,
		status,
		created_at,
		updated_at 
	FROM products 
	ORDER BY id DESC 
	%s %s`, lim, off)

	stmt, err = rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return 0, nil, err
	}

	rows, err := stmt.QueryContext(ctx)
	defer rows.Close()
	if err != nil {
		return 0, nil, err
	}
	for rows.Next() {
		var prod admin.Product

		err := rows.Scan(
			&prod.ID,
			&prod.UomID,
			&prod.CategoryID,
			&prod.Name,
			&prod.Description,
			&prod.DimensionDescription,
			&prod.Stock,
			&prod.Price,
			&prod.Discount,
			&prod.FinalPrice,
			&prod.IsHaveExpiry,
			&prod.ExpiredAt,
			&prod.ProductImage,
			&prod.Status,
			&prod.CreatedAt,
			&prod.UpdatedAt,
		)
		if err != nil {
			return 0, nil, err
		}
		prods = append(prods, &prod)
	}

	if len(prods) == 0 && counter.count > 0 {
		return counter.count, nil, errors.New("offset page")
	}

	return counter.count, prods, nil
}

func (rep *Repository) getCategory(ID int) (*admin.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r := map[string]int{"id": ID}
	paramRequest, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBufferString(string(paramRequest))
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/user/category", payload)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var category *admin.Category
	err = json.Unmarshal(body, &category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (rep *Repository) getCategories() ([]*admin.Category, error) {
	// prepare context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// make request config
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/user/categories", nil)
	if err != nil {
		return nil, err
	}

	// set request header
	req.Header.Add("Content-Type", "application/json")

	// make http client
	client := &http.Client{}

	// send request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var categories []*admin.Category
	err = json.Unmarshal(body, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
