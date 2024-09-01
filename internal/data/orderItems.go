package data

import "database/sql"

type OrderItem struct {
	ID        int64 `json:"id"`
	VariantID int64 `json:"variant_id"`
	Price     int32 `json:"price"`
	Quantity  int32 `json:"quantity"`
	OrderID   int64 `json:"order_id"`
	Version   int32 `json:"version"`
}

type OrderItemsModel struct {
	DB *sql.DB
}

func (m *OrderItemsModel) Insert(tx *sql.Tx, orderItem *OrderItem) error {

	query := `
		INSERT INTO order_items(variant_id, price, quantity, order_id)
		VALUES ($1, $2, $3, $4)
		returning id, version;
`
	args := []any{orderItem.VariantID, orderItem.Price, orderItem.Quantity, orderItem.OrderID}
	err := tx.QueryRow(query, args...).Scan(
		&orderItem.OrderID,
		&orderItem.Version)

	if err != nil {
		if err.Error() == "pq: insert or update on table \"order_items\" violates foreign key constraint \"order_items_order_id_fkey\"" {
			return ErrOrderNotFound
		}
		if err.Error() == "pq: insert or update on table \"order_items\" violates foreign key constraint \"order_items_variant_id_fkey\"" {
			return ErrVariantDoesNotExist
		}
		return err

	}
	return nil
}
