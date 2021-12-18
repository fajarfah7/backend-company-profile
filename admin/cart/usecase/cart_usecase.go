package usecase

import (
	"fmt"
	"net/http"
	"project_company_profile/models/admin"
	"project_company_profile/response"
)

type Usecase struct {
	Rep admin.CartRepository
}

func NewCartUsecase(rep admin.CartRepository) admin.CartUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) Cart(ID int64) (res response.Response) {
	cart, err := ucs.Rep.Cart(ID)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = cart
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "OK"}}
	return
}

func (ucs *Usecase) Carts(status, page, limit int) (pagination admin.CartPagination) {
	count, carts, err := ucs.Rep.Carts(status, page, limit)
	totalPage := 0
	prevPage := 0
	nextPage := 0
	if err != nil {
		pagination = admin.CartPagination{
			Count:        int64(0),
			PreviousPage: prevPage,
			CurrentPage:  page,
			NextPage:     nextPage,
			TotalPage:    totalPage,
			Data:         nil,
		}
		return
	}

	restPage := count % int64(limit)
	totalPage = int((count - int64(restPage)) / int64(limit))
	if restPage > 0 {
		totalPage++
	}

	if page-1 > 0 {
		prevPage = page - 1
	}

	if page+1 <= totalPage {
		nextPage = page + 1
	}

	pagination = admin.CartPagination{
		Count:        count,
		PreviousPage: prevPage,
		CurrentPage:  page,
		NextPage:     nextPage,
		TotalPage:    totalPage,
		Data:         carts,
	}
	return
}

func (ucs *Usecase) Update(ID int64, status int, shipmentCode string) (res response.Response) {
	fmt.Println("ucs", shipmentCode)
	err := ucs.Rep.Update(ID, status, shipmentCode)
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
