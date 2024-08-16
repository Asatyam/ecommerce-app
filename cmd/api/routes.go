package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.demo)
	router.HandlerFunc(http.MethodPost, "/user", app.createUserHandler)
	return router
}
