package usecase

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"project_company_profile/models/user"
	"project_company_profile/response"
	"strconv"
)

type Usecase struct {
	Rep user.CartRepository
}

func NewCartUsecase(rep user.CartRepository) user.CartUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) Create(userID int64) (res response.Response) {
	cart, err := ucs.Rep.Create(userID)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		return
	}
	res.ActionStatus = true
	res.Data = cart
	res.Messages = []map[string]string{{"api": "OK"}}
	res.Status = http.StatusOK
	return
}

func (ucs *Usecase) Update(ID int64, status int) (res response.Response) {
	paymentCode := encodeID(ID)

	cart, err := ucs.Rep.Update(ID, status, paymentCode)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		return
	}
	res.ActionStatus = true
	res.Data = cart
	res.Messages = []map[string]string{{"api": "OK"}}
	res.Status = http.StatusOK
	return
}

func (ucs *Usecase) GetActiveCart(userID int64) (res response.Response) {
	// err on here indicate that there is no active cart
	cart, err := ucs.Rep.GetActiveCart(userID)
	if err != nil {
		// create new cart
		createCart, errCreateCart := ucs.Rep.Create(userID)
		if errCreateCart != nil {
			res.ActionStatus = false
			res.Data = nil
			res.Messages = []map[string]string{{"api": err.Error()}}
			res.Status = http.StatusInternalServerError
			return
		}
		cart = createCart
	}
	res.ActionStatus = true
	res.Data = cart
	res.Messages = []map[string]string{{"api": "OK"}}
	res.Status = http.StatusOK
	return
}

func (ucs *Usecase) GetByID(cartID int64) (cart *user.Cart, status int) {
	resCart, err := ucs.Rep.GetByID(cartID)
	if err != nil {
		cart = nil
		status = http.StatusInternalServerError
		log.Println(err)
		return
	}
	status = http.StatusOK
	cart = resCart
	return
}

func (ucs *Usecase) Carts(userID int64, status, page, limit int, key string) (pagination user.CartPagination) {
	count, carts, err := ucs.Rep.Carts(userID, status, page, limit, key)
	totalPage := 0
	prevPage := 0
	nextPage := 0
	if err != nil {
		pagination = user.CartPagination{
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

	pagination = user.CartPagination{
		Count:        count,
		PreviousPage: prevPage,
		CurrentPage:  page,
		NextPage:     nextPage,
		TotalPage:    totalPage,
		Data:         carts,
	}
	return
}

func (ucs *Usecase) Search(userID int64, status int, key string) (cart []*user.Cart, httpStatusCode int) {
	carts, err := ucs.Rep.Search(userID, status, key)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return carts, http.StatusOK
}

func encodeID(ID int64) string {
	intID := int(ID)
	h := md5.New()
	io.WriteString(h, strconv.Itoa(intID))
	bytesID := h.Sum(nil)
	stringID := fmt.Sprintf("%x", bytesID)
	return stringID
}
