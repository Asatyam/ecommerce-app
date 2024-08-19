package data

import (
	"database/sql"
	"errors"
)

type Models struct {
	Users      UserModel
	Tokens     TokenModel
	Products   ProductModel
	Brands     BrandModel
	Categories CategoryModel
}

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

func NewModels(db *sql.DB) Models {
	return Models{
		Users:      UserModel{DB: db},
		Tokens:     TokenModel{DB: db},
		Products:   ProductModel{DB: db},
		Brands:     BrandModel{DB: db},
		Categories: CategoryModel{DB: db},
	}
}
