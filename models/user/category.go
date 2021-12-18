package user

import (
	"github.com/gin-gonic/gin"
)

type Category struct {
	ID   int    `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}

type CategoryDelivery interface {
	Category(c *gin.Context)
	Categories(c *gin.Context)
}

type CategoryUsecase interface {
	Category(ID int) (category *Category, httpStatus int)
	Categories() (categories []*Category, httpStatus int)
}

type CategoryRepository interface {
	Category(ID int) (category *Category, err error)
	Categories() (categories []*Category, err error)
}
