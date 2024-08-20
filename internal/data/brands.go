package data

import (
	"database/sql"
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/validator"
)

type Brand struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Logo        string `json:"logo"`
	Version     string `json:"-"`
}

type BrandModel struct {
	DB *sql.DB
}

var (
	ErrDuplicateBrand      = errors.New("duplicate brand")
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

func ValidateBrand(v *validator.Validator, brand *Brand) {

	v.Check(brand.Name != "", "name", "cannot be empty")
	v.Check(len(brand.Name) <= 50, "name", "must be less than 50 characters")
	v.Check(brand.Description != "", "description", "cannot be empty")
	v.Check(brand.Logo != "", "logo", "cannot be empty")

}

func (m *BrandModel) Insert(brand *Brand) error {

	query := `
		INSERT INTO brands(name, description, logo)  VALUES 
		($1, $2, $3)
		RETURNING id, version
`
	args := []any{brand.Name, brand.Description, brand.Logo}

	err := m.DB.QueryRow(query, args...).Scan(&brand.ID, &brand.Version)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "brands_name_key"` {
			return ErrDuplicateBrand
		}
		return err
	}
	return nil

}
func (m *BrandModel) Get(id int64) (*Brand, error) {

	query := `
			SELECT id, name, description, logo, version
			FROM brands
			WHERE id=$1
		`
	var brand Brand
	err := m.DB.QueryRow(query, id).Scan(
		&brand.ID,
		&brand.Name,
		&brand.Description,
		&brand.Logo,
		&brand.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &brand, nil
}

func (m *BrandModel) Update(brand *Brand) error {

	query := `
			UPDATE brands
			SET name=$1, description=$2, logo=$3, version = version + 1
			WHERE id=$4 and version = $5
			RETURNING version
`
	args := []any{brand.Name, brand.Description, brand.Logo, brand.ID, brand.Version}

	err := m.DB.QueryRow(query, args...).Scan(&brand.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "brands_name_key"` {
			return ErrDuplicateBrand
		}
		return err
	}
	return nil
}

func (m *BrandModel) Delete(id int64) error {
	query := `
		DELETE FROM brands
		WHERE id = $1
`
	row, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
