package admin

import (
	"project_company_profile/response"

	"time"

	"github.com/gin-gonic/gin"
)

type Cart struct {
	ID           int64      `json:"id"`
	UserID       int64      `json:"user_id"`
	Status       int        `json:"status"`
	PaymentCode  *string    `json:"payment_code"`
	ShipmentCode *string    `json:"shipment_code"`
	ReceivedAt   *time.Time `json:"received_at"`
}

type CartPagination struct {
	Count        int64   `json:"count"`
	PreviousPage int     `json:"previous_page"`
	CurrentPage  int     `json:"current_page"`
	NextPage     int     `json:"next_page"`
	TotalPage    int     `json:"total_page"`
	Data         []*Cart `json:"data"`
}

type CartDelivery interface {
	Cart(c *gin.Context)
	Carts(c *gin.Context)
	Update(c *gin.Context)
}
type CartUsecase interface {
	Cart(ID int64) (res response.Response)
	Carts(status, page, limit int) (pagination CartPagination)
	Update(ID int64, status int, shipmentCode string) (res response.Response)
}
type CartRepository interface {
	Cart(ID int64) (cart *Cart, err error)
	Carts(status, page, limit int) (totalData int64, carts []*Cart, err error)
	Update(ID int64, status int, shipmentCode string) (err error)
}
