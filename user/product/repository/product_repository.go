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
	"project_company_profile/models/user"
	"strings"
	"time"
)

type Repository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) user.ProductRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) GetByID(id int64) (*user.Product, error) {
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
	WHERE id = $1 AND status = true`

	var prod user.Product
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}
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
	} else {
		category, err := rep.getCategory(prod.CategoryID)
		if err != nil {
			return &prod, nil
		}
		prod.Category = category
	}

	return &prod, nil
}

func (rep *Repository) GetAll(categoryID, page, limit int, key string) (count int64, products []*user.Product, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if key != "" {
		where = fmt.Sprintf(" AND (LOWER(name) LIKE '%%%s%%' OR LOWER(description) LIKE '%%%s%%') ", key, key)
	}
	if categoryID > 0 {
		where += fmt.Sprintf(" AND category_id=%d ", categoryID)
	}

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
	qCount := `SELECT COUNT(id) FROM products WHERE status = true` + where
	stmt, err := rep.DB.PrepareContext(ctx, qCount)
	defer stmt.Close()

	if err != nil {
		return 0, nil, err
	}
	countRows := stmt.QueryRowContext(ctx)

	countRows.Scan(&counter.count)
	if counter.count == 0 {
		return 0, nil, errors.New("There is no data")
	}

	// statement to do query
	var prods []*user.Product
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
	WHERE status = true %s
	ORDER BY id DESC 
	%s %s`, where, lim, off)

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
		var prod user.Product

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
		return int64(counter.count), nil, errors.New("offset page")
	}

	if len(prods) > 0 {
		categories, err := rep.getCategories()
		if err != nil {
			fmt.Println("user/product/repository/ThreeNewest", err)
			return int64(counter.count), prods, nil
		}

		for keyProd, prod := range prods {
			for _, cat := range categories {
				if prod.CategoryID == cat.ID {
					prods[keyProd].Category = cat
				}
			}
		}

	}

	return int64(counter.count), prods, nil
}

func (rep *Repository) ThreeNewest() ([]*user.Product, error) {
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
	FROM products ORDER BY id DESC LIMIT 3`

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

	var prods []*user.Product
	for rows.Next() {
		var prod user.Product
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
			return nil, err
		}

		prods = append(prods, &prod)
	}

	if len(prods) > 0 {
		categories, err := rep.getCategories()
		if err != nil {
			fmt.Println("user/product/repository/ThreeNewest", err)
			return prods, nil
		}

		for keyProd, prod := range prods {
			for _, cat := range categories {
				if prod.CategoryID == cat.ID {
					prods[keyProd].Category = cat
				}
			}
		}

	}

	return prods, nil

}

func (rep *Repository) GetWhereIDIn(IDs []int64) (products []*user.Product, err error) {
	// open db connection
	// db, err := database.OpenDB()
	// defer db.Close()
	// if err != nil {
	// 	return nil, err
	// }
	// rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// statement to do query
	whereIn := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(IDs)), ", "), "[]")
	q := fmt.Sprintf(`SELECT 
		products.id,
		products.uom_id,
		products.name,
		products.description,
		products.dimension_description,
		products.stock,
		products.price,
		products.discount,
		products.final_price,
		products.is_have_expiry,
		products.expired_at,
		products.product_image,
		products.status,
		products.created_at,
		products.updated_at 
	FROM products WHERE id IN(%s)`, whereIn)

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

	for rows.Next() {
		var prod user.Product

		err := rows.Scan(
			&prod.ID,
			&prod.UomID,
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
			return nil, err
		}

		products = append(products, &prod)
	}

	return products, nil
}

func (rep *Repository) getCategory(ID int) (*user.Category, error) {
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

	if res.StatusCode != http.StatusOK {
		fmt.Println("404 not found")
		return nil, errors.New("404 not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var category *user.Category
	err = json.Unmarshal(body, &category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (rep *Repository) getCategories() ([]*user.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/categories", nil)
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

	if res.StatusCode != http.StatusOK {
		fmt.Println("404 not found")
		return nil, errors.New("404 not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var categories []*user.Category
	err = json.Unmarshal(body, &categories)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
