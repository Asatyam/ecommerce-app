package main

import (
	"errors"
	"fmt"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"net/http"
	"strconv"
)

func (app *application) createProductVariantHandler(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(32 << 30)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var variant data.ProductVariant
	variant.Variants = make(map[string]any, 10)

	for key, values := range r.Form {
		value := values[0]

		switch key {
		case "product_id":
			productID, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}
			variant.ProductID = int64(productID)
		case "price":
			price, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}
			variant.Price = int32(price)
		case "discount":
			discount, err := strconv.ParseFloat(value, 32)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}
			variant.Discount = float32(discount)
		case "quantity":
			quantity, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}
			variant.Quantity = int32(quantity)
		case "sku":
			variant.SKU = value
		default:
			variant.Variants[key] = value
		}

	}

	v := validator.New()
	data.ValidateVariant(v, &variant)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	imgURL, err := app.getImageURL(r, "image")
	if err != nil {
		if errors.Is(err, data.ErrUnsupportedFileType) {
			app.unsupportedMediaTypeResponse(w, r)
			return
		}
		if err.Error() == "error Retrieving File from the form" {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	variant.Variants["image"] = imgURL

	err = app.models.ProductVariants.Insert(&variant)
	if err != nil {
		if errors.Is(err, data.ErrProductDoesNotExist) {
			app.errorResponse(w, r, http.StatusNotFound, "product does not exist")
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("variants/%d", variant.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"variant": variant}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) showProductVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	variant, err := app.models.ProductVariants.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrVariantDoesNotExist) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"variant": variant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) deleteProductVariantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.models.ProductVariants.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "variant deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
