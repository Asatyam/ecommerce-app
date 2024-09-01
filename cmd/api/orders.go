package main

import (
	"errors"
	"fmt"
	"github.com/Asatyam/ecommerce-app/internal/data"
	"net/http"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {

	tx, err := app.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		}
	}()

	var input struct {
		PaymentStatus string           `json:"payment_status"`
		Total         int32            `json:"total"`
		ContactNo     string           `json:"contact_no"`
		Address       string           `json:"address"`
		OrderItems    []data.OrderItem `json:"order_items"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)

	order := &data.Order{
		PaymentStatus: input.PaymentStatus,
		Total:         0,
		ContactNo:     input.ContactNo,
		Address:       input.Address,
		CustomerID:    user.ID,
	}

	err = app.models.Orders.Insert(tx, order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	var total int32
	for _, item := range input.OrderItems {
		orderItem := &data.OrderItem{
			Price:     item.Price,
			Quantity:  item.Quantity,
			VariantID: item.VariantID,
			OrderID:   order.ID,
		}
		total += item.Price * item.Quantity

		err = app.models.OrderItems.Insert(tx, orderItem)
		if err != nil {
			if errors.Is(err, data.ErrOrderNotFound) {
				app.errorResponse(w, r, http.StatusNotFound, err)
				return
			}
			if errors.Is(err, data.ErrVariantDoesNotExist) {
				app.errorResponse(w, r, http.StatusNotFound, err)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	order.Total = total
	err = app.models.Orders.Update(tx, order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/orders/%d", order.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"order": order}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) showOrderHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	order, err := app.models.Orders.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrOrderNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	user := app.contextGetUser(r)
	if order.CustomerID != user.ID {
		app.notPermittedResponse(w, r)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {

	tx, err := app.models.Orders.DB.Begin()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		} else {
			err = tx.Commit()
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
		}
	}()

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order, err := app.models.Orders.Get(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		PaymentStatus *string `json:"payment_status"`
		ContactNo     *string `json:"contact_no"`
		Address       *string `json:"address"`
		Status        *string `json:"status"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.PaymentStatus != nil {
		order.PaymentStatus = *input.PaymentStatus
	}
	if input.ContactNo != nil {
		order.ContactNo = *input.ContactNo
	}
	if input.Address != nil {
		order.Address = *input.Address
	}
	if input.Status != nil {
		order.Status = *input.Status
	}

	err = app.models.Orders.Update(tx, order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"order": order}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
func (app *application) getAllOrdersHandler(w http.ResponseWriter, r *http.Request) {

	user := app.contextGetUser(r)
	orders, err := app.models.Orders.GetAll(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
