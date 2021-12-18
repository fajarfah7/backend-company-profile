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
	"project_company_profile/models/admin"
	"time"
	// "time"
)

type Repository struct {
	DB *sql.DB
}

func NewCartProductRepository(db *sql.DB) admin.CartProductsRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) CartProducts(cartID int64) ([]*admin.CartProduct, error) {
	cartProducts, err := rep.getCartProducts(cartID)
	if err != nil {
		return nil, err
	}

	listProductID := []int64{}
	for _, cartProduct := range cartProducts {
		listProductID = append(listProductID, int64(cartProduct.ProductID))
	}

	products, err := rep.getProducts(listProductID)
	if err != nil {
		return nil, err
	}

	for i, cartProduct := range cartProducts {
		for _, product := range products {
			if cartProduct.ProductID == product.ID {
				cartProducts[i].Product = product
			}
		}
	}

	return cartProducts, nil
}

func (rep *Repository) getCartProducts(cartID int64) ([]*admin.CartProduct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := `SELECT id, cart_id, product_id, amount FROM cart_products WHERE cart_id=$1 ORDER BY id`
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, cartID)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	var cartProducts []*admin.CartProduct
	for rows.Next() {
		var cartProduct admin.CartProduct
		err := rows.Scan(
			&cartProduct.ID,
			&cartProduct.CartID,
			&cartProduct.ProductID,
			&cartProduct.Amount,
		)
		if err != nil {
			return nil, err
		}

		cartProducts = append(cartProducts, &cartProduct)
	}

	return cartProducts, nil
}

func (rep *Repository) getProducts(productIDs []int64) ([]*admin.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if len(productIDs) == 0 {
		return nil, errors.New("There is no product")
	}
	prodIDs := map[string][]int64{"id": productIDs}

	reqParamsJSON, err := json.Marshal(prodIDs)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBufferString(string(reqParamsJSON))

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/products/where-id-in", payload)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("resp", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("body", err)
		return nil, err
	}

	var products []*admin.Product
	err = json.Unmarshal(body, &products)
	if err != nil {
		fmt.Println("Unmarshal", err)
		return nil, err
	}

	return products, nil
}
