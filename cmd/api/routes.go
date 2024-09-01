package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.requireActivated(app.demo))

	//User Routes
	router.HandlerFunc(http.MethodPost, "/user", app.createUserHandler)
	router.HandlerFunc(http.MethodPut, "/user/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/user/authenticate", app.authenticateTokenHandler)
	router.HandlerFunc(http.MethodPost, "/user/forget-password", app.forgetPasswordTokenHandler)
	router.HandlerFunc(http.MethodPut, "/user/password", app.resetPasswordHandler)

	// Product Routes
	router.HandlerFunc(http.MethodPost, "/products", app.createProductHandler)
	router.HandlerFunc(http.MethodGet, "/products/:id", app.showProductHandler)
	router.HandlerFunc(http.MethodPatch, "/products/:id", app.updateProductHandler)
	router.HandlerFunc(http.MethodDelete, "/products/:id", app.deleteProductHandler)
	router.HandlerFunc(http.MethodGet, "/products/:id/variants", app.getAllProductVariants)

	// Brand Routes
	router.HandlerFunc(http.MethodGet, "/brands/:id", app.showBrandHandler)
	router.HandlerFunc(http.MethodPost, "/brands", app.createBrandHandler)
	router.HandlerFunc(http.MethodPatch, "/brands/:id", app.updateBrandHandler)
	router.HandlerFunc(http.MethodDelete, "/brands/:id", app.deleteBrandHandler)

	// Category Routes
	router.HandlerFunc(http.MethodGet, "/categories/:id", app.showCategoryHandler)
	router.HandlerFunc(http.MethodPost, "/categories", app.createCategoryHandler)
	router.HandlerFunc(http.MethodPatch, "/categories/:id", app.updateCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/categories/:id", app.deleteCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/categories/:id/ancestors", app.getCategoryWithAncestorsHandler)

	// Product Variants Routes
	router.HandlerFunc(http.MethodPost, "/variants/", app.createProductVariantHandler)
	router.HandlerFunc(http.MethodGet, "/variants/:id", app.showProductVariantHandler)
	router.HandlerFunc(http.MethodDelete, "/variants/:id", app.deleteProductVariantHandler)

	//Order Routes
	router.HandlerFunc(http.MethodPost, "/orders", app.requireActivated(app.createOrderHandler))
	
	return app.authenticate(router)
}
