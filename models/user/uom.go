package user

import (
	"github.com/gin-gonic/gin"
)

type Uom struct {
	ID     int     `json:"id" form:"id"`
	Name   string  `json:"name" form:"name"`
	Amount float32 `json:"amount" form:"amount"`
}

type UomDelivery interface {
	Uom(c *gin.Context)
	Uoms(c *gin.Context)
}
type UomUsecase interface {
	Uom(ID int) (uom *Uom, httpStatus int)
	Uoms() (uoms []*Uom, httpStatus int)
}
type UomRepository interface {
	Uom(ID int) (uom *Uom, err error)
	Uoms() (uoms []*Uom, err error)
}
