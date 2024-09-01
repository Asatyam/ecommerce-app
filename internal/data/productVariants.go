package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/validator"
)

var (
	ErrProductDoesNotExist = errors.New("product does not  exist")
)

type ProductVariant struct {
	ID        int64          `json:"id"`
	ProductID int64          `json:"product_id"`
	Price     int32          `json:"price"`
	Discount  float32        `json:"discount"`
	SKU       string         `json:"sku"`
	Quantity  int32          `json:"quantity"`
	Variants  map[string]any `json:"variants"`
	Version   int32          `json:"version"`
}

type ProductVariantsModel struct {
	DB *sql.DB
}

func ValidateVariant(v *validator.Validator, variant *ProductVariant) {

	v.Check(variant.ProductID > 0, "product_id", "product_id is not valid")
	v.Check(variant.Price > 0, "price", "price is not valid")
	v.Check(variant.Discount >= 0, "discount", "discount is not valid")
	v.Check(variant.SKU != "", "sku", "sku cannot be empty")
	v.Check(variant.Quantity > 0, "quantity", "quantity is not valid")
	v.Check(len(variant.Variants) != 0, "variant", "variants cannot be empty")
}

func (m *ProductVariantsModel) Insert(variant *ProductVariant) error {

	variantsJson, err := json.Marshal(variant.Variants)
	if err != nil {
		return err
	}

	query := `
		INSERT into product_variants (product_id, price, discount, sku, quantity, variants )
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, version;	
`
	args := []any{variant.ProductID, variant.Price, variant.Discount, variant.SKU, variant.Quantity, variantsJson}

	err = m.DB.QueryRow(query, args...).Scan(&variant.ID, &variant.Version)
	if err != nil {
		if err.Error() == "pq: insert or update on table \"product_variants\" violates foreign key constraint \"product_variants_product_id_fkey\"" {
			return ErrProductDoesNotExist
		}
		return err
	}
	return nil
}
