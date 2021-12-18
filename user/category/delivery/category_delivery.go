package delivery

import (
	"project_company_profile/models/user"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs user.CategoryUsecase
}

func NewCategoryDelivery(ucs user.CategoryUsecase) user.CategoryDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.Engine, dlv user.CategoryDelivery) {
	r.POST("/category", dlv.Category)
	r.POST("/categories", dlv.Categories)
}

func (dlv *Delivery) Category(c *gin.Context) {
	type binderCategory struct {
		ID int `json:"id" form:"id"`
	}
	binder := new(binderCategory)
	c.Bind(binder)
	res, status := dlv.Ucs.Category(binder.ID)
	c.JSON(status, res)
}

func (dlv *Delivery) Categories(c *gin.Context) {
	res, status := dlv.Ucs.Categories()
	c.JSON(status, res)
}
