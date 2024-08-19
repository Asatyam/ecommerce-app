package main

import (
	"errors"
	"fmt"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"net/http"
)

func (app *application) createBrandHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	imgURL, err := app.getImageURL(r, "logo")
	if err != nil {
		if errors.Is(err, data.ErrUnsupportedFileType) {
			app.unsupportedMediaTypeResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	brand := &data.Brand{
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Logo:        imgURL,
	}

	v := validator.New()
	if data.ValidateBrand(v, brand); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Brands.Insert(brand)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateBrand) {
			v.AddError("name", "company already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/brands/%d", brand.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"brand": brand}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *application) showBrandHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	brand, err := app.models.Brands.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"brand": brand}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) updateBrandHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	brand, err := app.models.Brands.Get(id)

	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Logo        *string `json:"logo"`
	}
	err = app.readJSON(w, r, &input)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		brand.Name = *input.Name
	}
	if input.Description != nil {
		brand.Description = *input.Description
	}
	if input.Logo != nil {
		brand.Logo = *input.Logo
	}
	v := validator.New()
	if data.ValidateBrand(v, brand); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Brands.Update(brand)
	if err != nil {
		if errors.Is(err, data.ErrEditConflict) {
			app.editConflictResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"brand": brand}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
