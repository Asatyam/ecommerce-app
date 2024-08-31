package main

import (
	"errors"
	"fmt"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"net/http"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		CategoryID  int64  `json:"category_id"`
		BrandID     int64  `json:"brand_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	product := &data.Product{
		Name:        input.Name,
		Description: input.Description,
		CategoryID:  input.CategoryID,
		BrandID:     input.BrandID,
	}

	v := validator.New()
	data.ValidateProduct(v, product)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Refactor
	_, err = app.models.Categories.Get(input.CategoryID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, "category not found")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	_, err = app.models.Brands.Get(input.BrandID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, "brand not found")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Products.Insert(product)
	if err != nil {
		if errors.Is(err, data.ErrProductAlreadyExists) {
			v.AddError("name", "product already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/products/%d", product.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"product": product}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) showProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}
	product, err := app.models.Products.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	categories, err := app.models.Categories.GetWithAncestors(product.CategoryID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			panic("category not found")
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	brand, err := app.models.Brands.Get(product.BrandID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			panic("brand not found")
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"product": product, "categories": categories, "brand": brand}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	product, err := app.models.Products.Get(id)
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
		CategoryID  *int64  `json:"category_id"`
		BrandID     *int64  `json:"brand_id"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}
	if input.BrandID != nil {
		product.BrandID = *input.BrandID
	}
	v := validator.New()
	data.ValidateProduct(v, product)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Categories.Get(product.CategoryID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, "category not found")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	_, err = app.models.Brands.Get(product.BrandID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.errorResponse(w, r, http.StatusNotFound, "brand not found")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	fmt.Printf("product: %+v\n", product)
	err = app.models.Products.Update(product)
	if err != nil {
		if errors.Is(err, data.ErrEditConflict) {
			app.editConflictResponse(w, r)
			return
		}
		if errors.Is(err, data.ErrProductAlreadyExists) {
			v.AddError("name", "product already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"product": product}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Products.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
		}
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "product deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
