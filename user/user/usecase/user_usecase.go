package usecase

import (
	"fmt"
	"net/http"
	"project_company_profile/models/user"
	"project_company_profile/response"

	"golang.org/x/crypto/bcrypt"
)

type Usecase struct {
	Rep user.UserRepository
}

func NewUserUsecase(rep user.UserRepository) user.UserUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) Create(user user.User) (res response.Response) {
	err := ucs.Rep.Create(user)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = nil
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "Success create user"}}
	return
}

func (ucs *Usecase) Update(user user.User) (res response.Response) {
	err := ucs.Rep.Update(user)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = nil
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "Success update user"}}
	return
}

func (ucs *Usecase) UpdatePassword(userID int64, oldPassword, newPassword string) (res response.Response) {
	hashedPassword, err := ucs.Rep.CheckPassword(userID)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(*hashedPassword), []byte(oldPassword))
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}

	// parse old password
	newPasswordHashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": "Failed encrypt new password"}}
		return
	}

	err = ucs.Rep.UpdatePassword(userID, string(newPasswordHashed))
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = nil
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "OK"}}
	return
}

func (ucs *Usecase) Delete(userID int64) (res response.Response) {
	err := ucs.Rep.Delete(userID)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = nil
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "Success delete user"}}
	return
}

func (ucs *Usecase) GetByUsername(username string) (res response.Response) {
	user, err := ucs.Rep.GetByUsername(username)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = user
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "Got the user"}}
	return
}

func (ucs *Usecase) GetByUsernameAndPassword(username, password string) (user *user.User, status int, message string) {
	user, err := ucs.Rep.GetByUsername(username)

	if err != nil {
		status = http.StatusInternalServerError
		return nil, status, err.Error()
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// this error indicate that the password does not match
		fmt.Println(err.Error())
		status = http.StatusBadRequest
		return nil, status, "Wrong username or password"
	}

	status = http.StatusOK
	return user, status, "OK"
}
