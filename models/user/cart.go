package user

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
	// Create(c *gin.Context)
	Update(c *gin.Context)
	GetActiveCart(c *gin.Context)
	GetByID(c *gin.Context)
	Carts(c *gin.Context)
	Search(c *gin.Context)
}
type CartUsecase interface {
	Create(userID int64) (res response.Response)
	Update(ID int64, status int) (res response.Response)
	GetActiveCart(userID int64) (res response.Response)
	GetByID(cartID int64) (cart *Cart, httpStatus int)
	Carts(userID int64, status, page, limit int, key string) (pagination CartPagination)
	Search(userID int64, status int, key string) (cart []*Cart, httpStatusCode int)
}
type CartRepository interface {
	Create(userID int64) (*Cart, error)
	Update(ID int64, status int, paymentCode string) (*Cart, error)
	GetActiveCart(userID int64) (*Cart, error)
	GetByID(ID int64) (*Cart, error)
	Carts(userID int64, status, page, limit int, key string) (totalData int64, carts []*Cart, err error)
	Search(userID int64, status int, key string) (cart []*Cart, err error)
}
