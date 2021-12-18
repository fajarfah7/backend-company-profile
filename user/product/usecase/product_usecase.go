package usecase

import (
	"fmt"
	"net/http"
	"project_company_profile/models/user"
	"project_company_profile/response"
	"strings"
)

type Usecase struct {
	Rep user.ProductRepository
}

func NewProductUsecase(rep user.ProductRepository) user.ProductUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) GetByID(id int64) (res response.Response) {
	prod, err := ucs.Rep.GetByID(id)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"error": err.Error()}}
		res.Status = http.StatusInternalServerError
		return
	}

	res.ActionStatus = true
	res.Data = prod
	res.Messages = []map[string]string{{"server": "OK"}}
	res.Status = 200
	return
}

func (ucs *Usecase) GetAll(categoryID, page, limit int, key string) (pag user.ProductPagination) {
	count, prods, err := ucs.Rep.GetAll(categoryID, page, limit, strings.ToLower(key))
	fmt.Println(err)
	if err != nil {
		pag.Count = 0
		pag.PreviousPage = 0
		pag.CurrentPage = 0
		pag.NextPage = 0
		pag.TotalPage = 0
		pag.Data = nil
		return
	}

	// calculation to determine total pages
	remainingData := count % int64(limit)
	evenData := count - remainingData
	totalPage := int(evenData / int64(limit))
	if remainingData > 0 {
		totalPage++
	}

	// dermining previous page(if needed on frontend)
	previousPage := 0
	if page > 1 {
		previousPage = page - 1
	}

	// dermining next page(if needed on frontend)
	nextPage := 0
	if page < totalPage {
		nextPage = page + 1
	}

	pag.Count = count
	pag.PreviousPage = previousPage
	pag.CurrentPage = page
	pag.NextPage = nextPage
	pag.TotalPage = totalPage
	pag.Data = prods
	return
}

func (ucs *Usecase) ThreeNewest() (res response.Response) {
	prods, err := ucs.Rep.ThreeNewest()
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
		Data:         prods,
		Status:       http.StatusOK,
		Messages:     []map[string]string{{"api": "OK"}},
	}
	return
}

func (ucs *Usecase) GetWhereIDIn(IDs []int64) []*user.Product {
	prods, err := ucs.Rep.GetWhereIDIn(IDs)
	if err != nil {
		return nil
	}
	return prods
}
