package usecase

import (
	"net/http"
	"project_company_profile/models/admin"
	"project_company_profile/response"
)

type Usecase struct {
	Rep admin.CategoryRepository
}

func NewCategoryUsecase(rep admin.CategoryRepository) admin.CategoryUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) Create(name string) (res response.Response) {
	err := ucs.Rep.Create(name)
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
		Data:         nil,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}
	return
}

func (ucs *Usecase) Update(ID int, name string) (res response.Response) {
	err := ucs.Rep.Update(ID, name)
	if err != nil {
		res = response.Response{
			ActionStatus: false,
			Data:         nil,
			Status:       http.StatusInternalServerError,
			Messages:     []map[string]string{{"api": err.Error()}},
		}
	}
	res = response.Response{
		ActionStatus: true,
		Data:         nil,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}
	return
}

func (ucs *Usecase) Delete(ID int) (res response.Response) {
	err := ucs.Rep.Delete(ID)
	if err != nil {
		res = response.Response{
			ActionStatus: false,
			Data:         nil,
			Status:       http.StatusInternalServerError,
			Messages:     []map[string]string{{"api": err.Error()}},
		}
	}
	res = response.Response{
		ActionStatus: true,
		Data:         nil,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}
	return
}
