package usecase

import (
	"net/http"
	"project_company_profile/models/user"
)

type Usecase struct {
	Rep user.CategoryRepository
}

func NewCategoryUsecase(rep user.CategoryRepository) user.CategoryUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) Category(ID int) (*user.Category, int) {
	category, err := ucs.Rep.Category(ID)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return category, http.StatusOK
}

func (ucs *Usecase) Categories() ([]*user.Category, int) {
	categories, err := ucs.Rep.Categories()
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return categories, http.StatusOK
}
