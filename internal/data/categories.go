package data

import (
	"database/sql"
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/validator"
)

var (
	ErrDuplicateCategory = errors.New("category already exists")
)

type Category struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ParentID int64  `json:"parent"`
	Version  int    `json:"version"`
}
type CategoryModel struct {
	DB *sql.DB
}

func ValidateCategory(v *validator.Validator, category *Category) {

	v.Check(category.Name != "", "name", "must be provided")
	v.Check(len(category.Name) <= 500, "name", "must not be more than 500 bytes")

}

func (m CategoryModel) Insert(category *Category) error {

	query := `
			INSERT INTO categories(name, parent)
			VALUES ($1, -1)
			RETURNING id		
`
	err := m.DB.QueryRow(query, category.Name).Scan(&category.ID)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "categories_name_key"` {
			return ErrDuplicateCategory
		}
		return err
	}
	category.ParentID = -1
	return nil
}
func (m CategoryModel) Get(id int64) (*Category, error) {
	query := `
		SELECT id, name, parent, version 
		FROM categories
		WHERE id = $1
`
	var category Category
	err := m.DB.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.ParentID,
		&category.Version)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &category, nil
}
func (m CategoryModel) Update(category *Category) error {

	query := `
		UPDATE categories
		SET name = $1, parent = $2, version = version + 1
		WHERE id = $3 and version = $4
		RETURNING version
`
	args := []any{category.Name, category.ParentID, category.ID, category.Version}
	err := m.DB.QueryRow(query, args...).Scan(&category.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "categories_name_key"` {
			return ErrDuplicateCategory
		}
		return err
	}
	return nil
}
func (m CategoryModel) Delete(id int64) error {
	query := `
	DELETE FROM categories
	WHERE id = $1
	
`
	rows, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
func (m CategoryModel) GetWithAncestors(id int64) ([]*Category, error) {

	query := `
        WITH RECURSIVE category_hierarchy AS (
            SELECT 
                id, 
                name, 
                parent
            FROM 
                categories
            WHERE 
                id = $1

            UNION ALL

            SELECT 
                c.id, 
                c.name, 
                c.parent
            FROM 
                categories c
            INNER JOIN 
                category_hierarchy ch 
            ON 
                c.id = ch.parent
            WHERE 
                ch.parent != -1
        )
        SELECT * FROM category_hierarchy;
    `
	rows, err := m.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)
	var categories []*Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.ParentID)
		if err != nil {
			return nil, err
		}
		categories = append(categories, &category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil

}
