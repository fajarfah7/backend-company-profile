package usecase

import (
	"net/http"
	"project_company_profile/models/user"
	"project_company_profile/response"
)

type Usecase struct {
	Rep user.TextRepository
}

func NewTextUsecase(rep user.TextRepository) user.TextUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) Texts(textType string) (res response.Response) {
	texts, err := ucs.Rep.Texts(textType)
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
		Data:         texts,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}
	return
}
