package admin

import (
	"project_company_profile/response"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID          int64      `json:"id" form:"id"`
	Name        string     `json:"name" form:"name"`
	Username    string     `json:"username" form:"username"`
	Email       string     `json:"email" form:"email"`
	Address     string     `json:"address" form:"address"`
	PhoneNumber string     `json:"phone_number" form:"phone_number"`
	Password    string     `json:"-" form:"-"`
	CreatedAt   *time.Time `json:"-" form:"-"`
	UpdatedAt   *time.Time `json:"-" form:"-"`
}

type UserDelivery interface {
	GetByID(c *gin.Context)
}

type UserUsecase interface {
	GetByID(ID int64) response.Response //*User
}

type UserRepository interface {
	GetByID(ID int64) (*User, error)
}
