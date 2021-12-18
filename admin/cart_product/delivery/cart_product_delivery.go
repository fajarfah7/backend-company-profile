package delivery

import (
	"project_company_profile/models/admin"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs admin.CartProductsUsecase
}

func NewCartProductDelivery(ucs admin.CartProductsUsecase) admin.CartProductsDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.RouterGroup, dlv admin.CartProductsDelivery) {
	r.POST("/cart-products", dlv.CartProducts)
}

func (dlv *Delivery) CartProducts(c *gin.Context) {
	type cartProductBinder struct {
		CartID int `json:"cart_id" form:"cart_id"`
	}

	var binder cartProductBinder
	c.Bind(&binder)

	res := dlv.Ucs.CartProducts(int64(binder.CartID))

	c.JSON(res.Status, res)
}
