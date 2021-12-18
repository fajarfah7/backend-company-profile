package usecase

import (
	"net/http"
	"project_company_profile/models/user"
)

type Usecase struct {
	Rep user.UomRepository
}

func NewUomUsecase(rep user.UomRepository) user.UomUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) Uom(ID int) (*user.Uom, int) {
	uom, err := ucs.Rep.Uom(ID)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return uom, http.StatusOK
}

func (ucs *Usecase) Uoms() ([]*user.Uom, int) {
	uoms, err := ucs.Rep.Uoms()
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return uoms, http.StatusOK
}
