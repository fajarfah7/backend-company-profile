package delivery

import (
	"project_company_profile/models/admin"
	// "time"

	"github.com/gin-gonic/gin"
)

// Delivery struct which will have methods CreateOrUpdate and GetAll
type Delivery struct {
	Ucs admin.AboutUsUsecase
}

// AboutUsBinder is struct to bind incoming data from frontend
type AboutUsBinder struct {
	Name string `form:"name" json:"name"`
	Type string `form:"type" json:"type"`
	Text string `form:"text" json:"text"`
}

// NewAboutUsDelivery return admin.AboutUsDelivery(interface)
func NewAboutUsDelivery(ucs admin.AboutUsUsecase) admin.AboutUsDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

// NewRoutes create routes for about-us for admin page
func NewRoutes(r *gin.RouterGroup, dlv admin.AboutUsDelivery) {
	r.POST("/about-us/save", dlv.CreateOrUpdate)
	r.POST("/about-us", dlv.GetAll)
}

// CreateOrUpdate catch incoming data from frontend and then will be thrown to usecase
func (dlv *Delivery) CreateOrUpdate(c *gin.Context) {
	var binder AboutUsBinder
	c.BindJSON(&binder)

	response := dlv.Ucs.CreateOrUpdate(binder.Name, binder.Type, binder.Text)

	c.JSON(response.Status, response)
}

// GetAll call GetAll method on usecase
func (dlv *Delivery) GetAll(c *gin.Context) {
	response := dlv.Ucs.GetAll()
	c.JSON(response.Status, response)
}
