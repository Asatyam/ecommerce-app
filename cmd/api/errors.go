package main

import (
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.PrintError(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {

	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "An unexpected error has occurred"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The request resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, err.Error())
}
func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (app *application) invalidPasswordTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid or missing password token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) authenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
func (app *application) inactiveAccountResponse(w http.ResponseWriter, r *http.Request) {
	message := "your account must be activated to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
func (app *application) notPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "your user account does not have necessary permissions to access this resource"
	app.errorResponse(w, r, http.StatusForbidden, message)
}
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again "
	app.errorResponse(w, r, http.StatusConflict, message)
}
func (app *application) unsupportedMediaTypeResponse(w http.ResponseWriter, r *http.Request) {
	message := "unsupported media type"
	app.errorResponse(w, r, http.StatusUnsupportedMediaType, message)
}
func (app *application) categoryParentNotFound(w http.ResponseWriter, r *http.Request) {
	message := "parent category not found"
	app.errorResponse(w, r, http.StatusBadRequest, message)
}
