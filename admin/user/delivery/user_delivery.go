package delivery

import (
	"project_company_profile/models/admin"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs admin.UserUsecase
}

func NewAdminUserDelivery(ucs admin.UserUsecase) admin.UserDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.RouterGroup, dlv admin.UserDelivery) {
	r.POST("/user/get-by-id", dlv.GetByID)
}

func (dlv *Delivery) GetByID(c *gin.Context) {
	type binderGetByID struct {
		ID int `json:"id" form:"id"`
	}

	binder := new(binderGetByID)
	c.Bind(binder)

	res := dlv.Ucs.GetByID(int64(binder.ID))
	c.JSON(res.Status, res)
}
