package admin

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

type CartProductsDelivery interface {
	CartProducts(c *gin.Context)
}
type CartProductsUsecase interface {
	CartProducts(cartID int64) (res response.Response)
}
type CartProductsRepository interface {
	CartProducts(cartID int64) (cartProducts []*CartProduct, err error)
}
