package order

import (
	"database/sql"

	"github.com/ardhptr21/ecomgo/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) CreateOrder(userId int, order types.CreateOrderPayload) (int, error) {
	res, err := s.db.Exec("INSERT INTO orders (userId, total, status, address) VALUES (?,?,?,?)",
		userId, order.Total, order.Status, order.Address)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *Store) CreateOrderItem(orderId int, item types.CreateOrderItemPayload) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES (?,?,?,?)",
		orderId, item.ProductID, item.Quantity, item.Price)

	return err
}
