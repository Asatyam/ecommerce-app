package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.requireActivated(app.demo))
	router.HandlerFunc(http.MethodPost, "/user", app.createUserHandler)
	router.HandlerFunc(http.MethodPut, "/user/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/user/authenticate", app.authenticateTokenHandler)
	router.HandlerFunc(http.MethodPost, "/user/forget-password", app.forgetPasswordTokenHandler)
	router.HandlerFunc(http.MethodPut, "/user/password", app.resetPasswordHandler)
	
	return app.authenticate(router)
}
