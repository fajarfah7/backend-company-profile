package user

import (
	"project_company_profile/response"
	"time"

	"github.com/gin-gonic/gin"
)

// CartProduct is a data which contain data from cart_products and products table
type CartProduct struct {
	ID        int64     `json:"id"`
	CartID    int64     `json:"cart_id"`
	ProductID int64     `json:"product_id"`
	Amount    int       `json:"amount"`
	Product   *Product  `json:"product" form:"product"`
	CreatedAt time.Time `json:"-" form:"-"`
	UpdatedAt time.Time `json:"-" form:"-"`
}

// CartProducts is data which contain Cart, CartProducts, and Products
type CartProducts struct {
	Cart         *Cart          `json:"cart" form:"cart"`
	CartProducts []*CartProduct `json:"cart_products" form:"cart_products"`
}

type CartProductsDelivery interface {
	CreateOrUpdate(c *gin.Context)
	// Delete(c *gin.Context)
	CartItem(c *gin.Context)
	CartItems(c *gin.Context)
}
type CartProductsUsecase interface {
	CreateOrUpdate(cartID, productID int64, amount int) (res response.Response)
	// Delete(cartID, productID int64) (res response.Response)
	CartItem(cartID int64, token string) (res response.Response)
	CartItems(userID int64) (res response.Response)
}
type CartProductsRepository interface {
	Create(cartID, productID int64, amount int) error
	Update(ID int64, amount int) error
	Delete(cartID, productID int64) error
	// used to determine the action is it will do create or update
	GetCartProduct(cartID, productID int64) (*CartProduct, error)
	// for api purpose
	CartItem(cartID int64, token string) (cartProducts *CartProducts, err error)
	CartItems(userID int64) (cartProducts []*CartProduct, err error)
}
