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
	// "time"
)

type Repository struct {
	DB *sql.DB
}

func NewCartProductRepository(db *sql.DB) user.CartProductsRepository {
	return &Repository{
		DB: db,
	}
}

func (rep *Repository) Create(cartID, productID int64, amount int) error {
	// db, err := database.OpenDB()
	// defer db.Close()
	// if err != nil {
	// 	return err
	// }

	// rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "INSERT INTO cart_products(cart_id, product_id, amount, created_at, updated_at) VALUES($1, $2, $3, $4, $5)"

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, cartID, productID, amount, time.Now(), time.Now())
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

func (rep *Repository) Update(ID int64, amount int) error {
	// db, err := database.OpenDB()
	// defer db.Close()
	// if err != nil {
	// 	return err
	// }

	// rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "UPDATE cart_products SET amount=$1, updated_at=$2 WHERE id=$3"

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, amount, time.Now(), ID)
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

func (rep *Repository) Delete(cartID, productID int64) error {
	// db, err := database.OpenDB()
	// defer db.Close()
	// if err != nil {
	// 	return err
	// }

	// rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "DELETE FROM cart_products WHERE cart_id=$1 AND product_id=$2"

	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, cartID, productID)
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

func (rep *Repository) GetCartProduct(cartID, productID int64) (*user.CartProduct, error) {
	// db, err := database.OpenDB()
	// defer db.Close()
	// if err != nil {
	// 	return nil, err
	// }

	// rep.DB = db

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, cart_id, product_id, amount FROM cart_products WHERE cart_id=$1 AND product_id=$2"
	stmt, err := rep.DB.PrepareContext(ctx, q)
	defer stmt.Close()
	if err != nil {
		return nil, err
	}

	var cartProduct user.CartProduct
	err = stmt.QueryRowContext(ctx, cartID, productID).Scan(
		&cartProduct.ID,
		&cartProduct.CartID,
		&cartProduct.ProductID,
		&cartProduct.Amount,
	)

	if err != nil {
		return nil, err
	}

	return &cartProduct, nil
}

func (rep *Repository) CartItem(cartID int64, token string) (*user.CartProducts, error) {
	// *Cart
	cart, err := rep.getSingleCart(cartID, token)
	if err != nil {
		return nil, err
	}

	// []*CartProduct
	// get multiple cart_product from cart_products table
	cartProducts, err := rep.getMultipleCartProduct(cartID)
	if err != nil {
		return nil, err
	}

	// []*int64
	// get product IDs based on given cart
	sliceCartID := []int64{cartID}
	productIDs, err := rep.getMultipleProductID(sliceCartID)
	if err != nil {
		return nil, err
	}
	// []*Product
	// shoot product API to get multiple product based on product id was gotten
	products, err := rep.getMultipleProduct(productIDs)
	if err != nil {
		return nil, err
	}

	for keyCartProduct, cartProduct := range cartProducts {
		for _, product := range products {
			if cartProduct.ProductID == product.ID && cartProduct.Product == nil {
				cartProducts[keyCartProduct].Product = product
			}
		}
	}

	var resCartProduct user.CartProducts
	resCartProduct.Cart = cart
	resCartProduct.CartProducts = cartProducts

	return &resCartProduct, nil
}

func (rep *Repository) CartItems(cartID int64) ([]*user.CartProduct, error) {

	cartItems, err := rep.getMultipleCartProduct(cartID)
	if err != nil {
		return nil, err
	}

	listItemID := []*int64{}
	for _, cartItem := range cartItems {
		var id int64
		id = int64(cartItem.ProductID)
		listItemID = append(listItemID, &id)
	}

	items, err := rep.getMultipleProduct(listItemID)
	if err != nil {
		return nil, err
	}

	for i, cartItem := range cartItems {
		for _, item := range items {
			if cartItem.ProductID == item.ID {
				cartItems[i].Product = item
			}
		}
	}

	return cartItems, nil
}

// get single cart from carts table by id
func (rep *Repository) getSingleCart(cartID int64, token string) (*user.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	paramMap := map[string]*int64{"id": &cartID}

	paramJSON, err := json.Marshal(paramMap)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBufferString(string(paramJSON))

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/user/cart/get-by-id", payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var cart user.Cart
	err = json.Unmarshal(body, &cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// get multiple cart_product from cart_products table
func (rep *Repository) getMultipleCartProduct(cartID int64) ([]*user.CartProduct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	q := "SELECT id, cart_id, product_id, amount FROM cart_products WHERE cart_id=$1 ORDER By updated_at DESC"
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

	var cartProducts []*user.CartProduct
	for rows.Next() {
		var cartProduct user.CartProduct
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

// FOR MULTIPLE CART----------------------------------------------

// get multiple cart from carts table by user_id
func (rep *Repository) getMultipleCart(userID int64, token string) ([]*user.Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	paramUserID := map[string]int64{"user_id": userID}

	paramUserIDJSON, err := json.Marshal(paramUserID)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewBufferString(string(paramUserIDJSON))

	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:4000/user/carts", payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var carts []*user.Cart
	err = json.Unmarshal(body, &carts)
	if err != nil {
		return nil, err
	}

	return carts, nil
}

// get multiple cart_product from cart_products table but in many cart_id
func (rep *Repository) getMultipleCartProducts(cartIDs []int64) ([]*user.CartProduct, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	whereIn := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cartIDs)), ","), "[]")
	q := fmt.Sprintf(`SELECT id, cart_id, product_id, amount FROM cart_products WHERE cart_id IN(%s)`, whereIn)

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

	var cartProducts []*user.CartProduct
	for rows.Next() {
		var cartProduct user.CartProduct
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

// REUSABLE FUNCTIONS----------------------------------------------

// get multiple product_id from cart_products table based on given cart_id
func (rep *Repository) getMultipleProductID(cartIDs []int64) ([]*int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	whereIn := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(cartIDs)), ", "), "[]")
	q := fmt.Sprintf(`SELECT product_id FROM cart_products WHERE cart_id IN (%s)`, whereIn)
	// q := `SELECT product_id FROM cart_products WHERE cart_id IN ($1)`
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

	var productIDs []*int64
	for rows.Next() {
		var productID int64

		err := rows.Scan(&productID)
		if err != nil {
			return nil, err
		}

		productIDs = append(productIDs, &productID)
	}

	return productIDs, nil
}

// get multiple product from products table
func (rep *Repository) getMultipleProduct(productIDs []*int64) ([]*user.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if len(productIDs) == 0 {
		return nil, errors.New("There is no product")
	}
	prodIDs := map[string][]*int64{"id": productIDs}

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
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var products []*user.Product
	err = json.Unmarshal(body, &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}
