package delivery

import (
	"fmt"
	"net/http"
	"project_company_profile/models/user"
	"project_company_profile/response"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Delivery struct {
	// User user.User
	// Binder UserBinder
	Ucs user.UserUsecase
}

type UserBinder struct {
	ID          int     `json:"id" form:"id"`
	Name        string  `json:"name" form:"name"`
	Username    string  `json:"username" form:"username"`
	Email       string  `json:"email" form:"email"`
	Address     string  `json:"address" form:"address"`
	PhoneNumber string  `json:"phone_number" form:"phone_number"`
	Password    *string `json:"password" form:"password"`
	RePassword  *string `json:"re_password" form:"re_password"`
}

func NewUserDelivery(ucs user.UserUsecase) user.UserDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

func NewRoutes(r *gin.Engine, dlv user.UserDelivery) {
	r.POST("/user/create", dlv.Create)
	r.POST("/user/get-by-username-and-password", dlv.GetByUsernameAndPassword)
}

func NewSecureRoutes(r *gin.RouterGroup, dlv user.UserDelivery) {
	r.POST("/update", dlv.Update)
	r.POST("/update-password", dlv.UpdatePassword)
	r.POST("/delete", dlv.Delete)
	r.POST("/get-by-username", dlv.GetByUsername)
}

func (dlv *Delivery) Create(c *gin.Context) {
	var res response.Response

	binder := new(UserBinder)

	c.Bind(binder)

	user, errMsgs := validateInput(binder)

	if len(errMsgs) > 0 {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusBadRequest
		res.Messages = errMsgs
		c.JSON(res.Status, res)
		return
	}

	encryptPass, err := bcrypt.GenerateFromPassword([]byte(*binder.Password), 10)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"api": err.Error()})

		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = errMsgs
		c.JSON(res.Status, res)
		return
	}
	user.Password = string(encryptPass)

	timeNow := time.Now()
	user.CreatedAt = &timeNow
	user.UpdatedAt = &timeNow

	res = dlv.Ucs.Create(user)
	c.JSON(res.Status, res)
	return
}

func (dlv *Delivery) Update(c *gin.Context) {
	var res response.Response

	binder := new(UserBinder)

	c.Bind(binder)

	user, errMsgs := validateInput(binder)

	if len(errMsgs) > 0 {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusBadRequest
		res.Messages = errMsgs
		c.JSON(res.Status, res)
		return
	}

	timeNow := time.Now()
	user.UpdatedAt = &timeNow

	res = dlv.Ucs.Update(user)

	c.JSON(res.Status, res)
	return
}

func (dlv *Delivery) UpdatePassword(c *gin.Context) {
	var res response.Response

	type binderUpdatePassword struct {
		ID            int    `json:"id" form:"id"`
		OldPassword   string `json:"old_password" form:"old_password"`
		NewPassword   string `json:"new_password" form:"new_password"`
		ReNewPassword string `json:"re_new_password" form:"re_new_password"`
	}

	binder := new(binderUpdatePassword)

	c.Bind(binder)

	if binder.OldPassword == "" || binder.NewPassword == "" || binder.ReNewPassword == "" {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusBadRequest
		res.Messages = []map[string]string{{"api": "Please fill New password and Re-New Password"}}
		c.JSON(res.Status, res)
		return
	}

	if binder.NewPassword != "" && binder.ReNewPassword != "" {
		if binder.NewPassword != binder.ReNewPassword {
			res.ActionStatus = false
			res.Data = nil
			res.Status = http.StatusBadRequest
			res.Messages = []map[string]string{{"api": "Password does not match"}}
			c.JSON(res.Status, res)
			return
		}
	}

	// // parse old password
	// oldPassword, err := bcrypt.GenerateFromPassword([]byte(binder.OldPassword), 10)
	// if err != nil {
	// 	res.ActionStatus = false
	// 	res.Data = nil
	// 	res.Status = http.StatusInternalServerError
	// 	res.Messages = []map[string]string{{"api": "Failed encrypt password"}}
	// 	c.JSON(res.Status, res)
	// 	return
	// }

	// // check the password
	// err = dlv.Ucs.CheckPassword(int64(binder.ID), string(oldPassword))
	// if err != nil {
	// 	res.ActionStatus = false
	// 	res.Data = nil
	// 	res.Status = http.StatusBadRequest
	// 	res.Messages = []map[string]string{{"api": "Password is wrong"}}
	// 	c.JSON(res.Status, res)
	// 	return
	// }

	// // try to parse for new password
	// encryptPass, err := bcrypt.GenerateFromPassword([]byte(binder.NewPassword), 10)
	// if err != nil {
	// 	res.ActionStatus = false
	// 	res.Data = nil
	// 	res.Status = http.StatusInternalServerError
	// 	res.Messages = []mfmt.Println(res)ap[string]string{{"api": "Failed encrypt password"}}
	// 	c.JSON(res.Status, res)
	// 	return
	// }
	// newPassword := string(encryptPass)

	res = dlv.Ucs.UpdatePassword(int64(binder.ID), binder.OldPassword, binder.NewPassword)

	c.JSON(res.Status, res)
}

func (dlv *Delivery) Delete(c *gin.Context) {
	var res response.Response
	type UserIDBinder struct {
		ID string `json:"id" form:"username"`
	}
	var userIDBinder UserIDBinder
	c.Bind(&userIDBinder)
	id, err := strconv.Atoi(userIDBinder.ID)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": "Invalid user ID"}}
		c.JSON(res.Status, res)
		return
	}
	userID := int64(id)
	res = dlv.Ucs.Delete(userID)
	c.JSON(res.Status, res)
	return
}

func (dlv *Delivery) GetByUsername(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	res := dlv.Ucs.GetByUsername(claims["username"].(string))
	c.JSON(res.Status, res)
	return
}

func (dlv *Delivery) GetByUsernameAndPassword(c *gin.Context) {
	type Credentials struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}

	var credentials Credentials

	c.Bind(&credentials)

	username := credentials.Username
	password := credentials.Password

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Username and password can not be emtpy",
		})
		return
	}

	user, status, message := dlv.Ucs.GetByUsernameAndPassword(username, password)
	if status != http.StatusOK {
		c.JSON(status, gin.H{
			"message": message,
		})
		return
	}
	c.JSON(status, user)
	return
}

func validateInput(binder *UserBinder) (user user.User, errMsgs []map[string]string) {
	// id, err := strconv.Atoi(dlv.Binder.ID)
	// if err != nil {
	// 	errMsgs = append(errMsgs, map[string]string{"api": "Invalid user ID"})
	// }
	user.ID = int64(binder.ID)

	if binder.Name == "" {
		errMsgs = append(errMsgs, map[string]string{"name": "Name can not be empty!"})
	}
	user.Name = binder.Name

	if binder.Username == "" {
		errMsgs = append(errMsgs, map[string]string{"username": "Username can not be empty!"})
	}
	user.Username = binder.Username

	if binder.Email == "" {
		errMsgs = append(errMsgs, map[string]string{"email": "Email can not be empty!"})
	}
	user.Email = binder.Email

	if binder.Address == "" {
		errMsgs = append(errMsgs, map[string]string{"address": "Address can not be empty!"})
	}
	user.Address = binder.Address

	if binder.PhoneNumber == "" {
		errMsgs = append(errMsgs, map[string]string{"phone_number": "Phone Number can not be empty!"})
	}
	user.PhoneNumber = binder.PhoneNumber

	// if action == create
	if user.ID == 0 {
		if binder.Password == nil {
			errMsgs = append(errMsgs, map[string]string{"password": "Password can not be empty!"})
		}
		if binder.RePassword == nil {
			errMsgs = append(errMsgs, map[string]string{"re_password": "Re-Password can not be empty!"})
		}

		if binder.Password != nil && binder.RePassword != nil {
			if *binder.Password == "" {
				errMsgs = append(errMsgs, map[string]string{"password": "Password can not be empty"})
			}
			if *binder.RePassword == "" {
				errMsgs = append(errMsgs, map[string]string{"re_password": "Re-Password can not be empty"})
			}
			if *binder.Password != *binder.RePassword {
				fmt.Println(binder.Password)
				fmt.Println(binder.RePassword)
				errMsgs = append(errMsgs, map[string]string{"password": "Password do not match!"})
			}
		}

	}

	// if action == edit and user change password
	if user.ID > 0 && (binder.Password != nil || binder.RePassword != nil) {
		if *binder.Password != *binder.RePassword {
			errMsgs = append(errMsgs, map[string]string{"password": "Password do not match!"})
		}
	}

	return user, errMsgs
}
