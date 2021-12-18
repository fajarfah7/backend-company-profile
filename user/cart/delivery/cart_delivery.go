package delivery

import (
	"net/http"
	"project_company_profile/models/user"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs        user.CartUsecase
	CartBinder CartBinder
	Cart       user.Cart
}

type CartBinder struct {
	ID     int `json:"id" form:"id"`
	UserID int `json:"user_id" form:"user_id"`
	Status int `json:"status" form:"status"`
}

func NewCartDelivery(ucs user.CartUsecase) user.CartDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.RouterGroup, dlv user.CartDelivery) {
	// r.POST("/cart/create", dlv.Create) // user_id
	r.POST("/cart/update", dlv.Update) // id, status
	r.POST("/cart", dlv.GetActiveCart) // user_id

	// these two path will be used for another services
	r.POST("/cart/get-by-id", dlv.GetByID) // id
	r.POST("/carts", dlv.Carts)            // user_id
	r.POST("/carts/search", dlv.Search)
}

func (dlv *Delivery) Create(c *gin.Context) {
	type CartBinder struct {
		UserID int `json:"user_id" form:"user_id"`
	}
	cartBinder := new(CartBinder)
	c.Bind(cartBinder)

	res := dlv.Ucs.Create(int64(cartBinder.UserID))
	c.JSON(res.Status, res)
}

func (dlv *Delivery) Update(c *gin.Context) {
	type CartBinder struct {
		ID     int `json:"id" form:"id"`
		Status int `json:"status" form:"status"`
	}
	cartBinder := new(CartBinder)
	c.Bind(cartBinder)
	res := dlv.Ucs.Update(int64(cartBinder.ID), cartBinder.Status)
	c.JSON(res.Status, res)
}

func (dlv *Delivery) GetActiveCart(c *gin.Context) {
	type CartBinder struct {
		UserID int `json:"user_id" form:"user_id"`
	}
	cartBinder := new(CartBinder)
	c.Bind(cartBinder)

	res := dlv.Ucs.GetActiveCart(int64(cartBinder.UserID))

	c.JSON(res.Status, res)
}

func (dlv *Delivery) GetByID(c *gin.Context) {
	type IDBinder struct {
		ID int `json:"id" form:"id"`
	}
	var iDBinder IDBinder
	c.Bind(&iDBinder)

	ID := int64(iDBinder.ID)

	cart, status := dlv.Ucs.GetByID(ID)

	c.JSON(status, cart)
}

func (dlv *Delivery) Carts(c *gin.Context) {
	type binderCarts struct {
		UserID int    `json:"user_id" form:"user_id"`
		Status int    `json:"status" form:"status"`
		Page   int    `json:"page" form:"page"`
		Limit  int    `json:"limit" form:"limit"`
		Key    string `json:"key" form:"key"`
	}
	binder := new(binderCarts)

	c.Bind(binder)

	res := dlv.Ucs.Carts(int64(binder.UserID), binder.Status, binder.Page, binder.Limit, binder.Key)
	if res.Count > 0 {
		c.JSON(http.StatusOK, res)
		return
	}

	errRes := []map[string]string{{"api": "There is no data"}}
	c.JSON(http.StatusInternalServerError, errRes)
	return
}

func (dlv *Delivery) Search(c *gin.Context) {
	type binderSearch struct {
		UserID int    `json:"user_id" form:"user_id"`
		Status int    `json:"status" form:"status"`
		Key    string `json:"key" form:"key"`
	}
	binder := new(binderSearch)
	c.Bind(binder)

	cart, status := dlv.Ucs.Search(int64(binder.UserID), binder.Status, binder.Key)
	c.JSON(status, cart)
}
