package usecase

import (
	"net/http"
	"project_company_profile/models/user"
	"project_company_profile/response"
)

type Usecase struct {
	Rep user.AboutUsRepository
}

func NewAboutUsUsecase(rep user.AboutUsRepository) user.AboutUsUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) GetAll() (res response.Response) {
	aboutUs, err := ucs.Rep.GetAll()
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
		Data:         aboutUs,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}
	return
}
