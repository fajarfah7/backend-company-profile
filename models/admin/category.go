package admin

import (
	"project_company_profile/response"

	"github.com/gin-gonic/gin"
)

type Category struct {
	ID   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type CategoryDelivery interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type CategoryUsecase interface {
	Create(name string) (res response.Response)
	Update(ID int, name string) (res response.Response)
	Delete(ID int) (res response.Response)
}

type CategoryRepository interface {
	Create(name string) error
	Update(ID int, name string) error
	Delete(ID int) error
}
