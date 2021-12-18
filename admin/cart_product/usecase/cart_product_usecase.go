package usecase

import (
	"net/http"
	"project_company_profile/models/admin"
	"project_company_profile/response"
)

type Usecase struct {
	Rep admin.CartProductsRepository
}

func NewCartProductUsecase(rep admin.CartProductsRepository) admin.CartProductsUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) CartProducts(cartID int64) (res response.Response) {
	cartProducts, err := ucs.Rep.CartProducts(cartID)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = cartProducts
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "OK"}}
	return
}
