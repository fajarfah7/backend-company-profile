package usecase

import (
	"net/http"
	"project_company_profile/models/user"
	"project_company_profile/response"
)

type Usecase struct {
	Rep user.CartProductsRepository
}

func NewCartProductUsecase(rep user.CartProductsRepository) user.CartProductsUsecase {
	return &Usecase{
		Rep: rep,
	}
}

func (ucs *Usecase) CreateOrUpdate(cartID, productID int64, amount int) (res response.Response) {
	cartProduct, _ := ucs.Rep.GetCartProduct(cartID, productID)

	// if cart-product is exist will update amount
	if cartProduct != nil {
		// if current amount + new amount == 0 will be deleted
		if cartProduct.Amount+amount == 0 {
			err := ucs.Rep.Delete(cartID, productID)
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
		err := ucs.Rep.Update(cartProduct.ID, cartProduct.Amount+amount)
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

	// if cart-product does not exist will create new cart-product
	err := ucs.Rep.Create(cartID, productID, amount)
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

// func (ucs *Usecase) Delete(cartID, productID int64) (res response.Response) {
// 	err := ucs.Rep.Delete(cartID, productID)
// 	if err != nil {
// 		res.ActionStatus = false
// 		res.Data = nil
// 		res.Status = http.StatusInternalServerError
// 		res.Messages = []map[string]string{{"api": err.Error()}}
// 		return
// 	}
// 	res.ActionStatus = true
// 	res.Data = nil
// 	res.Status = http.StatusOK
// 	res.Messages = []map[string]string{{"api": "OK"}}
// 	return
// }
func (ucs *Usecase) CartItem(cartID int64, token string) (res response.Response) {
	cartProducts, err := ucs.Rep.CartItem(cartID, token)
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
func (ucs *Usecase) CartItems(cartID int64) (res response.Response) {
	cartItems, err := ucs.Rep.CartItems(cartID)
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Status = http.StatusInternalServerError
		res.Messages = []map[string]string{{"api": err.Error()}}
		return
	}
	res.ActionStatus = true
	res.Data = cartItems
	res.Status = http.StatusOK
	res.Messages = []map[string]string{{"api": "OK"}}
	return
}
