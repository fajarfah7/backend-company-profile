package delivery

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"project_company_profile/admin/product/repository"
	"project_company_profile/database"
	"project_company_profile/models/admin"
	"project_company_profile/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// pq: sorry, too many clients already

// Delivery struct which will have methods Create, Update, Delete, GetByID, and GetAll(this is for pagination purpose)
type Delivery struct {
	Prod       admin.Product
	ProdBinder ProductBinder
	Ucs        admin.ProductUsecase
}

// ProductBinder is used to bind data from the form
type ProductBinder struct {
	ID                   string `form:"id"`
	Name                 string `form:"name"`
	UomID                string `form:"uom_id"`
	CategoryID           string `form:"category_id"`
	Stock                string `form:"stock"`
	Price                string `form:"price"`
	Discount             string `form:"discount"`
	DimensionDescription string `form:"dimension_description"`
	IsHaveExpiry         bool   `form:"is_have_expiry"`
	ExpiredAt            string `form:"expired_at"`
	Description          string `form:"description"`
	ProductImagePath     string `form:"product_image_path"`
	Status               bool   `form:"status"`
	// ProductImage         *multipart.FileHeader `form:"product_image"`
}

// NewProductDelivery return admin.ProductDelivery(interface)
func NewProductDelivery(ucs admin.ProductUsecase) admin.ProductDelivery {
	return &Delivery{
		Ucs: ucs,
	}
}

// NewRoutes create routes for product for admin page
func NewRoutes(r *gin.RouterGroup, dlv admin.ProductDelivery) {
	r.POST("/product/create", dlv.Create)
	r.POST("/product/update", dlv.Update)
	r.POST("/product/delete", dlv.Delete)
	r.POST("/product/get-by-id", dlv.GetByID)
	r.POST("/products", dlv.GetAll)
}

// Create on here do validation and process the image
func (dlv *Delivery) Create(c *gin.Context) {
	var res response.Response
	// bind the data  to struct
	c.Bind(&dlv.ProdBinder)
	// validate input
	errMsgs := dlv.validateInput()

	if len(errMsgs) > 0 {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = errMsgs
		res.Status = http.StatusBadRequest
		c.JSON(res.Status, res)
		return
	}

	// processing image start from here
	// receive file input
	file, err := c.FormFile("product_image")
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"product_image": err.Error()})

		res.ActionStatus = false
		res.Data = nil
		res.Messages = errMsgs
		res.Status = http.StatusBadRequest
		c.JSON(res.Status, res)
		return
	}
	if file == nil {
		errMsgs = append(errMsgs, map[string]string{"product_image": "file is empty"})

		res.ActionStatus = false
		res.Data = nil
		res.Messages = errMsgs
		res.Status = http.StatusBadRequest
		c.JSON(res.Status, res)
		return
	}

	// set file name and path where the file will be saved
	// filepath.Base(file.Filename) will return filename on current path
	filename := filepath.Base(file.Filename)
	filepath := "./images/" + filename
	c.SaveUploadedFile(file, filepath)

	// checking if uploaded file is exist
	checkFile, err := os.Stat(filepath)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"api": err.Error()})

		res.ActionStatus = false
		res.Data = nil
		res.Messages = errMsgs
		res.Status = http.StatusInternalServerError
		c.JSON(res.Status, res)
		return
	}
	if checkFile == nil {
		errMsgs = append(errMsgs, map[string]string{"api": "failed put image"})

		res.ActionStatus = false
		res.Data = nil
		res.Messages = errMsgs
		res.Status = http.StatusInternalServerError
		c.JSON(res.Status, res)
		return
	}

	// set saved file type to the struct for response purpose (type string)
	dlv.Prod.ProductImage = filename
	// end process imge

	// set the usecase
	resp := dlv.Ucs.Create(dlv.Prod)

	// return the response
	c.JSON(http.StatusOK, resp)
}

// Update on here do validation and process the image
func (dlv *Delivery) Update(c *gin.Context) {
	var res response.Response

	c.Bind(&dlv.ProdBinder)
	// validate input
	errMsgs := dlv.validateInput()

	if len(errMsgs) > 0 {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = errMsgs
		res.Status = http.StatusBadRequest
		c.JSON(res.Status, res)
		return
	}

	// get old image name
	// // open databse
	// db, err := database.OpenDB()
	// if err != nil {
	// 	errMsgs = append(errMsgs, map[string]string{"api": err.Error()})

	// 	res.ActionStatus = false
	// 	res.Data = nil
	// 	res.Messages = errMsgs
	// 	res.Status = http.StatusInternalServerError
	// 	c.JSON(res.Status, res)
	// 	return
	// }
	// // create new repository variable
	// repo := repository.NewProductRepository(db)
	// // get product based on the id
	currentProd, err := dlv.Ucs.GetProductByID(dlv.Prod.ID)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"api": "failed get current data"})

		res.ActionStatus = false
		res.Data = nil
		res.Messages = errMsgs
		res.Status = http.StatusInternalServerError
		c.JSON(res.Status, res)
		return
	}
	// set product image
	dlv.Prod.ProductImage = currentProd.ProductImage
	//

	// processing image start from here
	// receive file input
	file, _ := c.FormFile("product_image")
	// if file is nil this means admin does not change product image

	if file != nil {
		// check old image
		oldImage, err := os.Stat("./images/" + dlv.Prod.ProductImage)
		if err != nil {
			errMsgs = append(errMsgs, map[string]string{"api": err.Error()})

			res.ActionStatus = false
			res.Data = nil
			res.Messages = errMsgs
			res.Status = http.StatusInternalServerError
			c.JSON(res.Status, res)
			return
		}
		if oldImage != nil {
			err := os.Remove("./images/" + dlv.Prod.ProductImage)
			if err != nil {
				errMsgs = append(errMsgs, map[string]string{"api": err.Error()})

				res.ActionStatus = false
				res.Data = nil
				res.Messages = errMsgs
				res.Status = http.StatusInternalServerError
				c.JSON(res.Status, res)
				return
			}
		}

		// set file name and path where the file will be saved
		// filepath.Base(file.Filename) will return filename on current path
		filename := filepath.Base(file.Filename)
		filepath := "./images/" + filename
		c.SaveUploadedFile(file, filepath)

		// checking if uploaded file is exist
		checkFile, err := os.Stat(filepath)
		if err != nil {
			errMsgs = append(errMsgs, map[string]string{"api": err.Error()})

			res.ActionStatus = false
			res.Data = nil
			res.Messages = errMsgs
			res.Status = http.StatusInternalServerError
			c.JSON(res.Status, res)
			return
		}
		if checkFile == nil {
			errMsgs = append(errMsgs, map[string]string{"api": "image is null"})

			res.ActionStatus = false
			res.Data = nil
			res.Messages = errMsgs
			res.Status = http.StatusInternalServerError
			c.JSON(res.Status, res)
			return
		}

		// set saved file type to the struct for response purpose (type string)
		dlv.Prod.ProductImage = filename
		// end process imge
	}

	// set the usecase
	resp := dlv.Ucs.Update(dlv.Prod)

	// return the response
	c.JSON(200, resp)
}

// Delete call Delete method on usecse, in this method also have inner struct
func (dlv *Delivery) Delete(c *gin.Context) {
	var res response.Response
	type PostID struct {
		ID string `json:"id" form:"id"`
	}
	var postID PostID
	c.Bind(&postID)
	id, err := strconv.Atoi(postID.ID)

	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		c.JSON(res.Status, res)
		return
	}

	db, err := database.OpenDB()
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		c.JSON(res.Status, res)
		return
	}

	repo := repository.NewProductRepository(db)
	prod, err := repo.GetByID(int64(id))
	if err != nil {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		c.JSON(res.Status, res)
		return
	}

	if prod.ProductImage == "" {
		resp := dlv.Ucs.Delete(int64(id))
		c.JSON(resp.Status, resp)
		return
	}

	checkProductImage, err := os.Stat("./images/" + prod.ProductImage)
	if errors.Is(err, os.ErrPermission) {
		res.ActionStatus = false
		res.Data = nil
		res.Messages = []map[string]string{{"api": err.Error()}}
		res.Status = http.StatusInternalServerError
		c.JSON(res.Status, res)
		return
	} else if errors.Is(err, os.ErrNotExist) {
		resp := dlv.Ucs.Delete(int64(id))
		c.JSON(resp.Status, resp)
		return
	}
	if checkProductImage != nil {
		err := os.Remove("./images/" + prod.ProductImage)
		if err != nil {
			res.ActionStatus = false
			res.Data = nil
			res.Messages = []map[string]string{{"api": err.Error()}}
			res.Status = http.StatusInternalServerError
			c.JSON(res.Status, res)
			return
		}
	}

	resp := dlv.Ucs.Delete(int64(id))
	c.JSON(resp.Status, resp)
	return
}

// GetByID catch the id from parameter from URL, method is GET
func (dlv *Delivery) GetByID(c *gin.Context) {
	type GetByIDBinder struct {
		ID int64 `json:"id" form:"id"`
	}
	var getByIDBinder GetByIDBinder
	c.Bind(&getByIDBinder)

	res := dlv.Ucs.GetByID(getByIDBinder.ID)
	c.JSON(res.Status, res)
}

// GetAll catch page and limit parameter from URL, method is GET
func (dlv *Delivery) GetAll(c *gin.Context) {
	type GetAllBinder struct {
		Page  int `json:"page" form:"page" `
		Limit int `json:"limit" form:"page"`
	}

	var getAllBinder GetAllBinder

	c.Bind(&getAllBinder)

	res := dlv.Ucs.GetAll(getAllBinder.Page, getAllBinder.Limit)
	c.JSON(http.StatusOK, res)
}

func (dlv *Delivery) validateInput() (errMsgs []map[string]string) {
	// id
	if dlv.ProdBinder.ID == "" {
		errMsgs = append(errMsgs, map[string]string{"id": "id is empty"})
	}
	id, err := strconv.Atoi(dlv.ProdBinder.ID)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"id": err.Error()})
	}
	dlv.Prod.ID = int64(id)

	categoryID, err := strconv.Atoi(dlv.ProdBinder.CategoryID)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"id": "Failed parse category ID"})
	}
	dlv.Prod.CategoryID = categoryID

	// uom_id
	if dlv.ProdBinder.UomID == "" {
		errMsgs = append(errMsgs, map[string]string{"uom_id": "uom_id is empty"})
	}
	uomID, err := strconv.Atoi(dlv.ProdBinder.UomID)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"uom_id": err.Error()})
	}
	dlv.Prod.UomID = uomID

	// name
	if dlv.ProdBinder.Name == "" {
		errMsgs = append(errMsgs, map[string]string{"name": "product name is empty"})
	}
	dlv.Prod.Name = dlv.ProdBinder.Name

	// description
	dlv.Prod.Description = dlv.ProdBinder.Description

	// dimension_description
	if dlv.ProdBinder.DimensionDescription == "" {
		errMsgs = append(errMsgs, map[string]string{"dimension_description": "product dimension description is empty"})
	}
	dlv.Prod.DimensionDescription = dlv.ProdBinder.DimensionDescription

	// stock
	if dlv.ProdBinder.Stock == "" {
		errMsgs = append(errMsgs, map[string]string{"stock": "product stock is empty"})
	}
	stock, err := strconv.Atoi(dlv.ProdBinder.Stock)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"stock": err.Error()})
	}
	dlv.Prod.Stock = int64(stock)

	// price
	if dlv.ProdBinder.Price == "" {
		errMsgs = append(errMsgs, map[string]string{"price": "product price is empty"})
	}
	price, err := strconv.ParseFloat(dlv.ProdBinder.Price, 2)
	if err != nil {
		errMsgs = append(errMsgs, map[string]string{"price": err.Error()})
	}
	dlv.Prod.Price = price

	// discount
	discount := 0
	if dlv.ProdBinder.Discount != "" {
		discountInt, err := strconv.Atoi(dlv.ProdBinder.Discount)
		if err != nil {
			errMsgs = append(errMsgs, map[string]string{"discount": err.Error()})
		} else {
			discount = discountInt
		}
	}
	dlv.Prod.Discount = discount

	// final_price will be managed on usecase

	dlv.Prod.FinalPrice = 0.00

	// is_have_expiry
	dlv.Prod.IsHaveExpiry = dlv.ProdBinder.IsHaveExpiry

	// expired_at
	var expiredAt *time.Time
	expiredAt = nil
	if dlv.Prod.IsHaveExpiry == true {
		parsedTime, err := time.Parse("2006-01-02", dlv.ProdBinder.ExpiredAt)
		if err != nil {
			errMsgs = append(errMsgs, map[string]string{"expired_at": err.Error()})
		} else {
			expiredAt = &parsedTime
		}
	}
	dlv.Prod.ExpiredAt = expiredAt

	// product_image
	// if dlv.Prod.ID == 0 {
	// 	if dlv.ProdBinder.ProductImage == nil {
	// 		errMsgs = append(errMsgs, map[string]string{"product_image": "product image is empty"})
	// 	}
	// 	dlv.Prod.ProductImage = dlv.ProdBinder.ProductImagePath
	// }

	// status
	dlv.Prod.Status = dlv.ProdBinder.Status

	// set current time
	timeNow := time.Now()
	// created_at
	if dlv.Prod.ID == 0 {
		dlv.Prod.CreatedAt = &timeNow
	} else {
		dlv.Prod.CreatedAt = nil
	}

	// updated_at
	dlv.Prod.UpdatedAt = &timeNow

	return errMsgs
}
