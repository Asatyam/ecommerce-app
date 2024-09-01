package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrOrderNotFound = errors.New("order not found")
)

type Order struct {
	ID            int64     `json:"id"`
	Status        string    `json:"status"`
	PaymentStatus string    `json:"payment_status"`
	Total         int32     `json:"total"`
	ContactNo     string    `json:"contact_no"`
	Date          time.Time `json:"date"`
	CustomerID    int64     `json:"customer_id"`
	Address       string    `json:"address"`
	Version       int32     `json:"version"`
}

type OrderModel struct {
	DB *sql.DB
}

func (m *OrderModel) Insert(tx *sql.Tx, order *Order) error {

	query := `
		INSERT INTO orders(payment_status, total, contact_no, customer_id, address)
		VALUES ($1, $2, $3, $4, $5)
		returning id, status, date, version;`

	args := []any{order.PaymentStatus, order.Total, order.ContactNo, order.CustomerID, order.Address}

	err := tx.QueryRow(query, args...).Scan(
		&order.ID,
		&order.Status,
		&order.Date,
		&order.Version)

	if err != nil {
		return err
	}
	return nil
}
func (m *OrderModel) Get(id int64) (*Order, error) {

	query := `
		SELECT id, status, payment_status, total, contact_no, customer_id, address, version
		FROM orders
		WHERE id = $1;
`
	var order Order
	err := m.DB.QueryRow(query, id).Scan(
		&order.ID,
		&order.Status,
		&order.PaymentStatus,
		&order.Total,
		&order.ContactNo,
		&order.CustomerID,
		&order.Address,
		&order.Version,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}

	return &order, nil

}
func (m *OrderModel) Update(tx *sql.Tx, order *Order) error {

	query := `
		UPDATE orders
		SET status = $1, payment_status = $2, total = $3, contact_no = $4, address = $5, version = version + 1
		WHERE id = $6
		returning version;
	`
	args := []any{order.Status, order.PaymentStatus, order.Total, order.ContactNo, order.Address, order.ID}
	err := tx.QueryRow(query, args...).Scan(
		&order.Version)

	if err != nil {
		return err
	}
	return nil
}
func (m *OrderModel) Delete(id int64) error {
	query := `
		DELETE FROM orders WHERE id = $1;
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
		return ErrOrderNotFound
	}
	return nil
}
func (m *OrderModel) GetAll(customerID int64) ([]*Order, error) {
	query := `
		SELECT id, status, payment_status, total, contact_no, customer_id, address, version
		FROM orders
		WHERE customer_id = $1;
`
	rows, err := m.DB.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			return
		}
	}(rows)

	var orders []*Order
	for rows.Next() {
		var order Order
		err = rows.Scan(
			&order.ID,
			&order.Status,
			&order.PaymentStatus,
			&order.Total,
			&order.ContactNo,
			&order.CustomerID,
			&order.Address,
			&order.Version)

		if err != nil {
			return nil, err
		}
		orders = append(orders, &order)

	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
