package usecase

import (
	"log"
	"net/http"
	"project_company_profile/models/admin"
	"project_company_profile/response"
)

// Usecase struct which will have methods Create, Update, Delete, GetByID and GetAll
type Usecase struct {
	Rep admin.ProductRepository
}

// NewProductUsecase return admin.NewProductUsecase interface
func NewProductUsecase(rep admin.ProductRepository) admin.ProductUsecase {
	return &Usecase{
		Rep: rep,
	}
}

// Create process and do some calculation before stored to database
func (ucs *Usecase) Create(prod admin.Product) (res response.Response) {
	if prod.Discount > 0 {
		finalPrice := float64(prod.Discount) / 100 * prod.Price
		prod.FinalPrice = prod.Price - finalPrice
	}

	err := ucs.Rep.Create(prod)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		return
	}

	res.ActionStatus = true
	res.Data = nil
	res.Messages = []map[string]string{{"api": "Success save data"}}
	res.Status = http.StatusOK
	return
}

// Update process and do some calculation, and will update existing data on database
func (ucs *Usecase) Update(prod admin.Product) (res response.Response) {
	if prod.Discount > 0 {
		finalPrice := float64(prod.Discount) / 100 * prod.Price
		prod.FinalPrice = prod.Price - finalPrice
	}

	err := ucs.Rep.Update(prod)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		return
	}

	res.ActionStatus = true
	res.Data = nil
	res.Messages = []map[string]string{{"api": "Success update data"}}
	res.Status = http.StatusOK
	return
}

// Delete call Delete method on repository to delete the data based on given ID
func (ucs *Usecase) Delete(id int64) (res response.Response) {
	err := ucs.Rep.Delete(id)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		return
	}

	res.ActionStatus = true
	res.Data = nil
	res.Messages = []map[string]string{{"api": "success delete data"}}
	res.Status = http.StatusOK
	return
}

// GetByID call method on repository to find data based on given ID
func (ucs *Usecase) GetByID(id int64) (res response.Response) {
	prod, err := ucs.Rep.GetByID(id)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		return
	}

	res.ActionStatus = true
	res.Data = prod
	res.Messages = []map[string]string{{"api": "OK"}}
	res.Status = http.StatusOK
	return
}

// GetProductByID call method on repository to find data based on given ID
func (ucs *Usecase) GetProductByID(id int64) (*admin.Product, error) {
	prod, err := ucs.Rep.GetProductByID(id)
	return prod, err
}

// GetAll call method on repository to get some data based on the limit(for pagination purpose)
func (ucs *Usecase) GetAll(page, limit int) (pag admin.ProductPagination) {
	count, prods, err := ucs.Rep.GetAll(page, limit)
	if err != nil {
		log.Println(err)
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
