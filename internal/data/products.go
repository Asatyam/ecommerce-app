package data

import (
	"database/sql"
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	BrandID     int64     `json:"brandId"`
	CategoryID  int64     `json:"categoryId"`
	Version     string    `json:"-"`
}

type ProductModel struct {
	DB *sql.DB
}

func ValidateProduct(v validator.Validator, product *Product) {

	v.Check(product.Name != "", "name", "is required")
	v.Check(len(product.Name) <= 100, "name", "size of name must not be greater than 100")
	v.Check(product.Description != "", "description", "is required")
	v.Check(product.BrandID > 0, "brand_id", "is required")
	v.Check(product.CategoryID > 0, "category_id", "is required")
}

func (m ProductModel) Insert(product Product) error {

	query := `
		INSERT INTO products (name, description, category_id, brand_id, version) VALUES 
		($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

	args := []any{product.Name, product.Description, product.CategoryID, product.BrandID, product.Version}
	err := m.DB.QueryRow(query, args...).Scan(
		&product.ID,
		&product.CreatedAt,
		&product.Version)

	return err
}
func (m ProductModel) Get(id int64) (*Product, error) {

	query := `
		SELECT id, name, description, category_id, brand_id, created_at 
		FROM products
		WHERE id = $1
		`
	var product Product
	err := m.DB.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&product.BrandID,
		&product.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &product, nil
}
func (m ProductModel) Update(product Product) error {

	query := `
		UPDATE products
		SET name = $1, description = $2, brand_id = $3, category_id = $4, version = version + 1
		WHERE id = $5 and version = $6
		RETURNING  version
`
	args := []any{product.Name, product.Description, product.BrandID, product.CategoryID, product.ID, product.Version}
	err := m.DB.QueryRow(query, args...).Scan(
		&product.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}
func (m ProductModel) Delete(id int64) error {
	query := `
		DELETE FROM products
		WHERE id = $1
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
