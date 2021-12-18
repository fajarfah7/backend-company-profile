package delivery

import (
	"project_company_profile/models/user"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs user.AboutUsUsecase
}

func NewAboutUsDelivery(ucs user.AboutUsUsecase) user.AboutUsDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.Engine, dlv user.AboutUsDelivery) {
	r.POST("/about-us", dlv.GetAll)
}

func (dlv *Delivery) GetAll(c *gin.Context) {
	res := dlv.Ucs.GetAll()
	c.JSON(res.Status, res)
}
