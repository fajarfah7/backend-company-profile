package user

import (
	"project_company_profile/response"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/lib/pq"
)

// To avoid nil value when call and scan the data
// 1) can use sql.NullString, then to call it AboutUs.Text.String
// 2) set the type to be pointer(*string) then to call it *AboutUs.Text
// 3) when select the data use coalesce, select id, coalesce(text, ''), so if the text is nil will return empty text

// AboutUs is a struct which has design same as table "texts" on database
type AboutUs struct {
	ID        int64      `json:"id"`
	ImageID   int64      `json:"image_id"` // to avoid nil value can use sql.NullInt64 type
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// AboutUsDelivery catch request and give response to frontend
type AboutUsDelivery interface {
	GetAll(*gin.Context)
}

// AboutUsUsecase will handle input and output data
type AboutUsUsecase interface {
	GetAll() response.Response
}

// AboutUsRepository handle communication with database
type AboutUsRepository interface {
	GetAll() ([]*AboutUs, error)
}
