package main

import (
	"errors"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"github.com/Asatyam/ecommerce-app/internal/validator"
	"net/http"
)

func (app *application) showCategoryHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	category, err := app.models.Categories.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	category := &data.Category{
		Name: input.Name,
	}
	v := validator.New()
	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Categories.Insert(category)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateCategory) {
			v.AddError("name", "category name already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.badRequestResponse(w, r, err)
		return
	}
	category, err := app.models.Categories.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	var input struct {
		Name     *string `json:"name"`
		ParentID *int64  `json:"parent"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.ParentID != nil {
		if *input.ParentID != -1 {
			_, err := app.models.Categories.Get(*input.ParentID)
			if err != nil {
				if errors.Is(err, data.ErrRecordNotFound) {
					app.categoryParentNotFound(w, r)
					return
				}
				app.serverErrorResponse(w, r, err)
				return
			}
		}
		category.ParentID = *input.ParentID
	}
	v := validator.New()
	if data.ValidateCategory(v, category); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	if category.ID == category.ParentID {
		v.AddError("parent_id", "must not be the same as its own id")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Categories.Update(category)
	if err != nil {
		if errors.Is(err, data.ErrEditConflict) {
			app.editConflictResponse(w, r)
			return
		}
		if errors.Is(err, data.ErrDuplicateCategory) {
			v.AddError("name", "category name already exists")
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.models.Categories.Delete(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "category has been deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *application) getCategoryWithAncestorsHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	categories, err := app.models.Categories.GetWithAncestors(id)
	if err != nil {
		if categories == nil {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"categories": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
