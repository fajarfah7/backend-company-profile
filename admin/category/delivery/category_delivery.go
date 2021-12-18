package delivery

import (
	"project_company_profile/models/admin"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs admin.CategoryUsecase
}

func NewCategoryDelivery(ucs admin.CategoryUsecase) admin.CategoryDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.RouterGroup, dlv admin.CategoryDelivery) {
	r.POST("/category/create", dlv.Create)
	r.POST("/category/update", dlv.Update)
	r.POST("/category/delete", dlv.Delete)
}

func (dlv *Delivery) Create(c *gin.Context) {
	type binderCreate struct {
		Name string `json:"name" form:"name"`
	}

	binder := new(binderCreate)
	c.Bind(binder)
	res := dlv.Ucs.Create(binder.Name)
	c.JSON(res.Status, res)
}

func (dlv *Delivery) Update(c *gin.Context) {
	type binderUpdate struct {
		ID   int    `json:"id" form:"id"`
		Name string `json:"name" form:"name"`
	}

	binder := new(binderUpdate)
	c.Bind(binder)
	res := dlv.Ucs.Update(binder.ID, binder.Name)
	c.JSON(res.Status, res)
}

func (dlv *Delivery) Delete(c *gin.Context) {
	type binderDelete struct {
		ID int `json:"id" form:"id"`
	}

	binder := new(binderDelete)
	c.Bind(binder)
	res := dlv.Ucs.Delete(binder.ID)
	c.JSON(res.Status, res)
}
