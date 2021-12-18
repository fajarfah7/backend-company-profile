package delivery

import (
	"net/http"
	"project_company_profile/models/admin"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs admin.CartUsecase
}

func NewCartDelivery(ucs admin.CartUsecase) admin.CartDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}
func NewRoutes(r *gin.RouterGroup, dlv admin.CartDelivery) {
	r.POST("/cart", dlv.Cart)
	r.POST("/carts", dlv.Carts)
	r.POST("/cart/update", dlv.Update)
}

func (dlv *Delivery) Cart(c *gin.Context) {
	type binderCart struct {
		ID int64 `json:"id" form:"id"`
	}
	binder := new(binderCart)
	c.Bind(binder)
	res := dlv.Ucs.Cart(binder.ID)
	c.JSON(res.Status, res)
	return
}

func (dlv *Delivery) Carts(c *gin.Context) {
	type binderCarts struct {
		Status int `json:"status" form:"status"`
		Page   int `json:"page" form:"page"`
		Limit  int `json:"limit" form:"limit"`
	}
	binder := new(binderCarts)

	c.Bind(binder)

	res := dlv.Ucs.Carts(binder.Status, binder.Page, binder.Limit)
	if res.Count > 0 {
		c.JSON(http.StatusOK, res)
		return
	}

	errRes := []map[string]string{{"api": "There is no data"}}
	c.JSON(http.StatusInternalServerError, errRes)
	return
}

func (dlv *Delivery) Update(c *gin.Context) {
	type binderCart struct {
		ID           int64  `json:"id" form:"id"`
		Status       int    `json:"Status" form:"status"`
		ShipmentCode string `json:"shipment_code" form:"shipment_code"`
	}
	binder := new(binderCart)
	c.Bind(binder)
	if binder.ShipmentCode == "" {
		messages := []map[string]string{{"api": "Shipment code can not be empty"}}
		c.JSON(http.StatusBadRequest, gin.H{
			"messages": messages,
		})
		return
	}
	res := dlv.Ucs.Update(binder.ID, binder.Status, binder.ShipmentCode)
	c.JSON(res.Status, res)
	return
}
