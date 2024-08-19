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

	return app.authenticate(router)
}
