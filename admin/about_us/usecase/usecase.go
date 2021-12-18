package usecase

import (
	"net/http"
	"project_company_profile/models/admin"
	"project_company_profile/response"
)

// Usecase struct which will have methods CreateOrUpdate and GetAll
type Usecase struct {
	Repo admin.AboutUsRepository
}

// NewAboutUsUsecase return admin.NewAboutUsUsecase interface
func NewAboutUsUsecase(rep admin.AboutUsRepository) admin.AboutUsUsecase {
	return &Usecase{
		Repo: rep,
	}
}

// CreateOrUpdate will call method CreateOrUpdate on repository
func (usc *Usecase) CreateOrUpdate(textName, textType, newText string) (res response.Response) {
	err := usc.Repo.CreateOrUpdate(textName, textType, newText)

	if err != nil {
		res.ActionStatus = false
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Data = nil
		return
	}

	res.ActionStatus = true
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "OK"}}
	res.Data = nil
	return
}

// GetAll will call method GetAll on repository
func (usc *Usecase) GetAll() (res response.Response) {
	data, err := usc.Repo.GetAll()
	if err != nil {
		res.ActionStatus = false
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Data = nil
		return
	}

	res.ActionStatus = true
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "OK"}}
	res.Data = data
	return
}
