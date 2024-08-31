package data

import (
	"database/sql"
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/validator"
)

var (
	ErrVariantAlreadyExists = errors.New("variant already exists")
)

type ProductVariant struct {
	ID        int64          `json:"id"`
	ProductID int64          `json:"product_id"`
	Price     float32        `json:"price"`
	Discount  float32        `json:"discount"`
	SKU       string         `json:"sku"`
	Variants  map[string]any `json:"variants"`
}

type ProductVariantsModel struct {
	DB *sql.DB
}

func ValidateVariant(v *validator.Validator, variant *ProductVariant) {
	v.Check(variant.ProductID > 0, "product_id", "product_id is not valid")
	v.Check(variant.Price > 0, "price", "price is not valid")
	v.Check(variant.Discount >= 0, "discount", "discount is not valid")
	v.Check(variant.SKU == "", "sku", "sku cannot be empty")
}
