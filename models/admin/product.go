package admin

import (
	"time"

	"project_company_profile/response"

	"github.com/gin-gonic/gin"
)

// ProductPagination response for pagination
type ProductPagination struct {
	Count        int64      `json:"count"`
	PreviousPage int        `json:"previous_page"`
	CurrentPage  int        `json:"current_page"`
	NextPage     int        `json:"next_page"`
	TotalPage    int        `json:"total_page"`
	Data         []*Product `json:"data"`
}

// Product is struct, table on DB is products
type Product struct {
	ID                   int64      `json:"id"`
	UomID                int        `json:"uom_id"`
	CategoryID           int        `json:"category_id"`
	UomName              string     `json:"uom_name"`
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

// ProductDelivery catch request and give response to frontend
type ProductDelivery interface {
	Create(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
	GetByID(*gin.Context)
	GetAll(*gin.Context)
}

// ProductUsecase process the data and then will be stored to DB
type ProductUsecase interface {
	Create(Product) response.Response
	Update(Product) response.Response
	Delete(int64) response.Response
	GetByID(id int64) response.Response
	GetProductByID(id int64) (*Product, error)
	GetAll(page, limit int) ProductPagination
}

// ProductRepository do communication to database
type ProductRepository interface {
	Create(Product) error
	Update(Product) error
	Delete(int64) error
	GetByID(id int64) (*Product, error)
	GetProductByID(id int64) (*Product, error)
	GetAll(page, limit int) (int64, []*Product, error)
}
