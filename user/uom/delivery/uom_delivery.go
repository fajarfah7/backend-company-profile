package delivery

import (
	"project_company_profile/models/user"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs user.UomUsecase
}

func NewUomDelivery(ucs user.UomUsecase) user.UomDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.Engine, dlv user.UomDelivery) {
	r.POST("/uom", dlv.Uom)
	r.POST("/uoms", dlv.Uoms)
}

func (dlv *Delivery) Uom(c *gin.Context) {
	type binderUom struct {
		ID int `json:"id" form:"id"`
	}
	binder := new(binderUom)

	c.Bind(binder)

	res, status := dlv.Ucs.Uom(binder.ID)
	c.JSON(status, res)
}

func (dlv *Delivery) Uoms(c *gin.Context) {
	res, status := dlv.Ucs.Uoms()
	c.JSON(status, res)
}
