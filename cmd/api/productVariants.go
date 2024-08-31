package main

import "net/http"

func (app *application) createProductVariant(w http.ResponseWriter, r *http.Request) {

	var input struct {
		ProductID string  `json:"product_id"`
		Price     float32 `json:"price"`
		Discount  float32 `json:"discount"`
		SKU       string  `json:"sku"`
	}
}
