package user

import (
	"project_company_profile/response"
	"time"

	"github.com/gin-gonic/gin"
)

type Text struct {
	ID        int64      `json:"id"`
	ImageID   int64      `json:"image_id"` // to avoid nil value can use sql.NullInt64 type
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Text      string     `json:"text"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// TextDelivery catch request and give response to frontend
type TextDelivery interface {
	Texts(c *gin.Context)
}

// TextUsecase will handle input and output data
type TextUsecase interface {
	Texts(textType string) (res response.Response)
}

// TextRepository handle communication with database
type TextRepository interface {
	Texts(textType string) (texts []*Text, err error)
}
