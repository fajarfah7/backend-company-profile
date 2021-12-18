package delivery

import (
	"project_company_profile/models/user"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs user.CartProductsUsecase
}

func NewCartProductDelivery(ucs user.CartProductsUsecase) user.CartProductsDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.RouterGroup, dlv user.CartProductsDelivery) {
	r.POST("/cart-product/save", dlv.CreateOrUpdate)
	// r.POST("/cart-product/delete", dlv.Delete)
	r.POST("/cart-item", dlv.CartItem)
	r.POST("/cart-items", dlv.CartItems)
}

// /cart-product/cart-products

// /cart-product/list-cart-products
type CartProductBinder struct {
	CartID    int `json:"cart_id" form:"cart_id"`
	ProductID int `json:"product_id" form:"product_id"`
	Amount    int `json:"amount" form:"amount"`
}

func (dlv *Delivery) CreateOrUpdate(c *gin.Context) {
	var cartProdBinder CartProductBinder

	c.Bind(&cartProdBinder)

	res := dlv.Ucs.CreateOrUpdate(int64(cartProdBinder.CartID), int64(cartProdBinder.ProductID), cartProdBinder.Amount)

	c.JSON(res.Status, res)
}

// func (dlv *Delivery) Delete(c *gin.Context) {
// 	var cartProdBinder CartProductBinder

// 	c.Bind(&cartProdBinder)

// 	res := dlv.Ucs.Delete(int64(cartProdBinder.CartID), int64(cartProdBinder.ProductID))

// 	c.JSON(res.Status, res)
// }

func (dlv *Delivery) CartItem(c *gin.Context) {
	type CartIDBinder struct {
		CartID int `json:"cart_id" form:"cart_id"`
	}

	var cartIDBinder CartIDBinder
	authorization := c.GetHeader("Authorization")
	c.Bind(&cartIDBinder)

	res := dlv.Ucs.CartItem(int64(cartIDBinder.CartID), authorization)

	c.JSON(res.Status, res)
}

func (dlv *Delivery) CartItems(c *gin.Context) {
	type binderCartItems struct {
		CartID int `json:"cart_id" form:"cart_id"`
	}

	binder := new(binderCartItems)
	c.Bind(binder)
	res := dlv.Ucs.CartItems(int64(binder.CartID))

	c.JSON(res.Status, res)
}
