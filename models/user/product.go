package user

import (
	"time"

	"project_company_profile/response"

	"github.com/gin-gonic/gin"
)

type ProductPagination struct {
	Count        int64      `json:"count"`
	PreviousPage int        `json:"previous_page"`
	CurrentPage  int        `json:"current_page"`
	NextPage     int        `json:"next_page"`
	TotalPage    int        `json:"total_page"`
	Data         []*Product `json:"data"`
}

type Product struct {
	ID                   int64      `json:"id"`
	UomID                int        `json:"uom_id"`
	CategoryID           int        `json:"category_id"`
	Category             *Category  `json:"category"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	DimensionDescription string     `json:"dimension_description"`
	Stock                int64      `json:"stock"`
	Price                float64    `json:"price"`
	Discount             int        `json:"discount"`
	FinalPrice           float64    `json:"final_price"`
	IsHaveExpiry         bool       `json:"is_have_expiry"`
	ExpiredAt            *time.Time `json:"expired_at"`
	ProductImage         string     `json:"product_image"`
	Status               bool       `json:"status"`
	CreatedAt            *time.Time `json:"-"`
	UpdatedAt            *time.Time `json:"-"`
}

type ProductDelivery interface {
	GetByID(c *gin.Context)
	GetAll(c *gin.Context)
	ThreeNewest(c *gin.Context)
	GetWhereIDIn(c *gin.Context)
}

type ProductUsecase interface {
	GetByID(id int64) (res response.Response)
	GetAll(categoryID, page, limit int, key string) (prodPagination ProductPagination)
	ThreeNewest() (res response.Response)
	// for another service
	GetWhereIDIn(IDs []int64) []*Product
}

type ProductRepository interface {
	GetByID(id int64) (*Product, error)
	GetAll(categoryID, page, limit int, key string) (count int64, products []*Product, err error)
	ThreeNewest() (threeNewest []*Product, err error)
	GetWhereIDIn(IDs []int64) (products []*Product, err error)
}
