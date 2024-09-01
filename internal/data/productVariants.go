package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/validator"
)

var (
	ErrProductDoesNotExist = errors.New("product does not  exist")
	ErrVariantDoesNotExist = errors.New("Variant does not  exist")
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

func (m *ProductVariantsModel) Get(id int64) (*ProductVariant, error) {

	query := `
	SELECT id, product_id, price, discount, sku, quantity, variants, version
	FROM product_variants
	WHERE id = $1;
	               	                                                                                                   
`
	var variant ProductVariant
	var jsonData []byte
	err := m.DB.QueryRow(query, id).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.Price,
		&variant.Discount,
		&variant.SKU,
		&variant.Quantity,
		&jsonData,
		&variant.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrVariantDoesNotExist
		}
		return nil, err
	}
	err = json.Unmarshal(jsonData, &variant.Variants)
	if err != nil {
		return nil, err
	}
	return &variant, nil
}

func (m *ProductVariantsModel) Delete(id int64) error {

	query := `
	DELETE FROM product_variants
	WHERE id = $1;
`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
