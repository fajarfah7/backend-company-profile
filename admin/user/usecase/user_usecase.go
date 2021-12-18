package usecase

import (
	"net/http"
	"project_company_profile/models/admin"
	"project_company_profile/response"
)

type Usecase struct {
	Rep admin.UserRepository
}

func NewAdminUserUsecase(rep admin.UserRepository) admin.UserUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) GetByID(ID int64) (res response.Response) {
	user, err := ucs.Rep.GetByID(ID)
	if err != nil {
		res = response.Response{
			ActionStatus: false,
			Data:         nil,
			Status:       http.StatusInternalServerError,
			Messages:     []map[string]string{{"api": err.Error()}},
		}
		return
	}

	res = response.Response{
		ActionStatus: true,
		Data:         user,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}
	return
}
