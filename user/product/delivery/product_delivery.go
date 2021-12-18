package delivery

import (
	"net/http"
	"project_company_profile/models/user"

	"github.com/gin-gonic/gin"
)

type Delivery struct {
	Ucs user.ProductUsecase
}

func NewProductDelivery(ucs user.ProductUsecase) user.ProductDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.Engine, dlv user.ProductDelivery) {
	r.POST("/product/get-by-id", dlv.GetByID)
	r.POST("/products", dlv.GetAll)
	r.POST("/product/three-newest", dlv.ThreeNewest)

	// used by another service
	r.POST("/products/where-id-in", dlv.GetWhereIDIn)
}

func (dlv *Delivery) GetByID(c *gin.Context) {
	type GetByIDBinder struct {
		ID int64 `json:"id" form:"id"`
	}
	var getByIDBinder GetByIDBinder
	c.Bind(&getByIDBinder)

	res := dlv.Ucs.GetByID(getByIDBinder.ID)
	c.JSON(res.Status, res)
}

func (dlv *Delivery) GetAll(c *gin.Context) {
	type binderGetAll struct {
		CategoryID int    `json:"category_id" form:"category_id"`
		Page       int    `json:"page" form:"page" `
		Limit      int    `json:"limit" form:"page"`
		Key        string `json:"key" form:"key"`
	}

	binder := new(binderGetAll)

	c.Bind(&binder)

	res := dlv.Ucs.GetAll(binder.CategoryID, binder.Page, binder.Limit, binder.Key)
	c.JSON(http.StatusOK, res)
}

func (dlv *Delivery) GetWhereIDIn(c *gin.Context) {
	type ListProductID struct {
		ID []int `json:"id"`
	}
	var listProdID ListProductID
	c.Bind(&listProdID)

	ids := []int64{}
	for _, val := range listProdID.ID {
		ids = append(ids, int64(val))
	}

	prods := dlv.Ucs.GetWhereIDIn(ids)

	c.JSON(http.StatusOK, prods)
}

func (dlv *Delivery) ThreeNewest(c *gin.Context) {
	res := dlv.Ucs.ThreeNewest()
	c.JSON(res.Status, res)
}
