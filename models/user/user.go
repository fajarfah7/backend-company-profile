package user

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
	Create(c *gin.Context)
	Update(c *gin.Context)
	UpdatePassword(c *gin.Context)
	Delete(c *gin.Context)
	GetByUsername(c *gin.Context)
	GetByUsernameAndPassword(c *gin.Context)
}

type UserUsecase interface {
	Create(user User) response.Response
	Update(user User) response.Response
	// CheckPassword(userID int64, oldPassword string) error
	UpdatePassword(userID int64, oldPassword, newPassword string) response.Response
	Delete(userID int64) response.Response
	GetByUsername(username string) response.Response //*User
	GetByUsernameAndPassword(username, password string) (user *User, httpStatus int, message string)
}

type UserRepository interface {
	Create(user User) error
	Update(user User) error
	CheckPassword(userID int64) (hashedPassword *string, err error)
	UpdatePassword(userID int64, newPassword string) error
	Delete(userID int64) error
	GetByUsername(username string) (*User, error)
	// GetByUsernameAndPassword(username, password string) (*User, error)
}
