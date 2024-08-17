package main

import (
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"github.com/asaskevich/govalidator"
	"net/http"
	"time"
)

func (app *application) authenticateTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	v.Check(input.Email != "", "email", "email must not be empty")
	v.Check(govalidator.IsEmail(input.Email), "email", "email is not valid")
	v.Check(input.Password != "", "password", "password must not be empty")
	v.Check(len(input.Password) >= 8, "password", "password must be of at least 8 characters")
	v.Check(len(input.Password) <= 500, "password", "password must be of at most 500 characters")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	ok, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !ok {
		app.invalidCredentialsResponse(w, r)
		return
	}
	token, err := app.models.Tokens.New(user.ID, data.ScopeAuthentication, 24*15*time.Hour)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *application) forgetPasswordTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	v.Check(input.Email != "", "email", "email must not be empty")
	v.Check(govalidator.IsEmail(input.Email), "email", "email is not valid")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.invalidCredentialsResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	if !user.Activated {
		app.inactiveAccountResponse(w, r)
		return
	}
	token, err := app.models.Tokens.New(user.ID, data.ScopePasswordReset, 10*time.Minute)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	data.ValidateToken(v, input.Token)
	v.Check(input.Password != "", "password", "password must not be empty")
	v.Check(len(input.Password) >= 8, "password", "password must be of at least 8 characters")
	v.Check(len(input.Password) <= 500, "password", "password must be of at most 500 characters")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	user, err := app.models.Users.GetForToken(data.ScopePasswordReset, input.Token)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.invalidPasswordTokenResponse(w, r)
			return
		}
		return
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Users.Update(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.models.Tokens.DeleteForUser(data.ScopePasswordReset, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Password was successfully reset!", "newPassword": input.Password}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
