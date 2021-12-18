package main

import (
	"log"
	"net/http"
	"project_company_profile/database"
	"project_company_profile/middlewares"
	_adminAuthMiddleware "project_company_profile/middlewares/admin"

	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/gin-gonic/gin"

	// admin
	// admin/category
	_adminCategoryDlv "project_company_profile/admin/category/delivery"
	_adminCategoryRep "project_company_profile/admin/category/repository"
	_adminCategoryUcs "project_company_profile/admin/category/usecase"

	// admin/product
	_adminProductDlv "project_company_profile/admin/product/delivery"
	_adminProductRep "project_company_profile/admin/product/repository"
	_adminProductUcs "project_company_profile/admin/product/usecase"

	_adminAboutUsDlv "project_company_profile/admin/about_us/delivery"
	_adminAboutUsRep "project_company_profile/admin/about_us/repository"
	_adminAboutUsUcs "project_company_profile/admin/about_us/usecase"

	// admin/cart
	_adminCartDlv "project_company_profile/admin/cart/delivery"
	_adminCartRep "project_company_profile/admin/cart/repository"
	_adminCartUcs "project_company_profile/admin/cart/usecase"

	// admin/cart-product
	_adminCartProductDlv "project_company_profile/admin/cart_product/delivery"
	_adminCartProductRep "project_company_profile/admin/cart_product/repository"
	_adminCartProductUcs "project_company_profile/admin/cart_product/usecase"

	// admin/user
	_adminUserDlv "project_company_profile/admin/user/delivery"
	_adminUserRep "project_company_profile/admin/user/repository"
	_adminUserUcs "project_company_profile/admin/user/usecase"

	// user
	// user/middleware/auth
	_userAuthMiddleware "project_company_profile/middlewares/user"

	//user/about-us
	_userAboutUsDlv "project_company_profile/user/about_us/delivery"
	_userAboutUsRep "project_company_profile/user/about_us/repository"
	_userAboutUsUcs "project_company_profile/user/about_us/usecase"

	// user/category
	_userCategoryDlv "project_company_profile/user/category/delivery"
	_userCategoryRep "project_company_profile/user/category/repository"
	_userCategoryUcs "project_company_profile/user/category/usecase"

	_userUserDlv "project_company_profile/user/user/delivery"
	_userUserRep "project_company_profile/user/user/repository"
	_userUserUcs "project_company_profile/user/user/usecase"

	// user/cart
	_userCartDlv "project_company_profile/user/cart/delivery"
	_userCartRep "project_company_profile/user/cart/repository"
	_userCartUcs "project_company_profile/user/cart/usecase"

	// user/product
	_userProductDlv "project_company_profile/user/product/delivery"
	_userProductRep "project_company_profile/user/product/repository"
	_userProductUcs "project_company_profile/user/product/usecase"

	// user/cart-products
	_userCartProductDlv "project_company_profile/user/cart_product/delivery"
	_userCartProductRep "project_company_profile/user/cart_product/repository"
	_userCartProductUcs "project_company_profile/user/cart_product/usecase"
)

func main() {
	// curl -X POST -H "Content-Type: application/json" -d '{"username":"admin", "password":"securepassword"}' http://localhost:4000/admin/login

	// valInt := []int{1,2,3,4}
	// valStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(valInt)), ", "), "[]")
	// fmt.Println(valStr)

	// CREATE/UPDATE/DELETE
	// res, err := ExecContext
	// affectedRow, err := res.RowsAffected()

	// GET SINGLE ROW
	// row := QueryRowContext
	// row.Scan(...)

	// GET MULTIPLE ROWS
	// rows, err := QueryContext
	// for rows.Next{...}

	// CHECK FOR IDLE CONNECTION ON POSTGRESQL
	// SELECT
	// datname,
	// pid,
	// application_name,
	// state,
	// query
	// FROM pg_stat_activity
	// WHERE datname='company_profile';

	var cors middlewares.CORS
	r := gin.Default()
	r.Use(cors.CORSMiddleware())

	// show images
	r.Static("/images", "./images")

	// open database
	db, err := database.OpenDB()
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// for admin routes

	// admin about us
	adminAboutUsRep := _adminAboutUsRep.NewAboutUsRepository(db)
	adminAboutUsUcs := _adminAboutUsUcs.NewAboutUsUsecase(adminAboutUsRep)
	adminAboutUsDlv := _adminAboutUsDlv.NewAboutUsDelivery(adminAboutUsUcs)

	// admin category
	adminCategoryRep := _adminCategoryRep.NewCategoryRepository(db)
	adminCategoryUcs := _adminCategoryUcs.NewCategoryUsecase(adminCategoryRep)
	adminCategoryDlv := _adminCategoryDlv.NewCategoryDelivery(adminCategoryUcs)

	// admin product
	adminProductRep := _adminProductRep.NewProductRepository(db)
	adminProductUcs := _adminProductUcs.NewProductUsecase(adminProductRep)
	adminProductDlv := _adminProductDlv.NewProductDelivery(adminProductUcs)

	// admin cart
	adminCartRep := _adminCartRep.NewCartRepository(db)
	adminCartUcs := _adminCartUcs.NewCartUsecase(adminCartRep)
	adminCartDlv := _adminCartDlv.NewCartDelivery(adminCartUcs)

	// admin cart products
	adminCartProductRep := _adminCartProductRep.NewCartProductRepository(db)
	adminCartProductUcs := _adminCartProductUcs.NewCartProductUsecase(adminCartProductRep)
	adminCartProductDlv := _adminCartProductDlv.NewCartProductDelivery(adminCartProductUcs)

	adminUserRep := _adminUserRep.NewAdminUserRepository(db)
	adminUserUcs := _adminUserUcs.NewAdminUserUsecase(adminUserRep)
	adminUserDlv := _adminUserDlv.NewAdminUserDelivery(adminUserUcs)

	// admin auth middleware for admin routes
	adminAuthMiddleware, err := _adminAuthMiddleware.AuthMiddleware()
	if err != nil {
		log.Println("JWT Error: " + err.Error())
	}

	errInitMiddlewareAdmin := adminAuthMiddleware.MiddlewareInit()
	if errInitMiddlewareAdmin != nil {
		log.Println("adminAuthMiddleware.MiddlewareInit() Error: " + errInitMiddlewareAdmin.Error())
	}

	r.POST("/admin/login", adminAuthMiddleware.LoginHandler)

	r.NoRoute(adminAuthMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	adminRoutes := r.Group("/admin")
	adminRoutes.GET("/refresh_token", adminAuthMiddleware.RefreshHandler)
	adminRoutes.Use(adminAuthMiddleware.MiddlewareFunc())
	{
		adminRoutes.POST("", func(c *gin.Context) {
			claims := jwt.ExtractClaims(c)
			c.JSON(http.StatusOK, gin.H{
				"id":   claims["id"],
				"name": claims["name"],
			})
			return
		})
		_adminCategoryDlv.NewRoutes(adminRoutes, adminCategoryDlv)
		_adminAboutUsDlv.NewRoutes(adminRoutes, adminAboutUsDlv)
		_adminProductDlv.NewRoutes(adminRoutes, adminProductDlv)
		_adminCartDlv.NewRoutes(adminRoutes, adminCartDlv)
		_adminCartProductDlv.NewRoutes(adminRoutes, adminCartProductDlv)
		_adminUserDlv.NewRoutes(adminRoutes, adminUserDlv)

	}

	// for user routes
	// user/user
	userUserRep := _userUserRep.NewUserRepository(db)
	userUserUcs := _userUserUcs.NewUserUsecase(userUserRep)
	userUserDlv := _userUserDlv.NewUserDelivery(userUserUcs)
	_userUserDlv.NewRoutes(r, userUserDlv)

	// user/about-us
	userAboutUsRep := _userAboutUsRep.NewAboutUsRepository(db)
	userAboutUsUcs := _userAboutUsUcs.NewAboutUsUsecase(userAboutUsRep)
	userAboutUsDlv := _userAboutUsDlv.NewAboutUsDelivery(userAboutUsUcs)
	_userAboutUsDlv.NewRoutes(r, userAboutUsDlv)

	// user/category
	userCategoryRep := _userCategoryRep.NewCategoryRepository(db)
	userCategoryUcs := _userCategoryUcs.NewCategoryUsecase(userCategoryRep)
	userCategoryDlv := _userCategoryDlv.NewCategoryDelivery(userCategoryUcs)
	_userCategoryDlv.NewRoutes(r, userCategoryDlv)

	// user/cart
	userCartRep := _userCartRep.NewCartRepository(db)
	userCartUcs := _userCartUcs.NewCartUsecase(userCartRep)
	userCartDlv := _userCartDlv.NewCartDelivery(userCartUcs)

	// user/cart-product
	userCartProductRep := _userCartProductRep.NewCartProductRepository(db)
	userCartProductUcs := _userCartProductUcs.NewCartProductUsecase(userCartProductRep)
	userCartProductDlv := _userCartProductDlv.NewCartProductDelivery(userCartProductUcs)

	// user/products
	userProductRep := _userProductRep.NewProductRepository(db)
	userProductUcs := _userProductUcs.NewProductUsecase(userProductRep)
	userProductDlv := _userProductDlv.NewProductDelivery(userProductUcs)
	_userProductDlv.NewRoutes(r, userProductDlv)

	userAuthMiddleware, err := _userAuthMiddleware.AuthMiddleware()
	if err != nil {
		log.Println("JWT Error: " + err.Error())
	}

	errInitMiddlewareUser := userAuthMiddleware.MiddlewareInit()
	if errInitMiddlewareUser != nil {
		log.Println("userAuthMiddleware.MiddlewareInit() Error: " + errInitMiddlewareUser.Error())
	}
	r.POST("/user/login", userAuthMiddleware.LoginHandler)
	r.NoRoute(userAuthMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	userRoutes := r.Group("/user")
	userRoutes.GET("/refresh_token", userAuthMiddleware.RefreshHandler)
	userRoutes.Use(userAuthMiddleware.MiddlewareFunc())
	{
		// userRoutes.POST("/init", func(c *gin.Context) {
		// 	claims := jwt.ExtractClaims(c)
		// 	c.JSON(http.StatusOK, claims)
		// 	return
		// })
		_userCartDlv.NewRoutes(userRoutes, userCartDlv)
		_userCartProductDlv.NewRoutes(userRoutes, userCartProductDlv)
		_userUserDlv.NewSecureRoutes(userRoutes, userUserDlv)
	}

	// defer db.Close()

	r.Run(":4000")
}
