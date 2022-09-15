package order

import (
	"database/sql"
	"errors"
)

type orderRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) GetOrdersByUser(userID string) ([]*OrderEntity, error) {
	return make([]*OrderEntity, 0), nil
}

func (o *orderRepository) Save(orderEntity *OrderEntity) error {
	return nil
}

func (o *orderRepository) FindByID(orderID string) (*OrderEntity, error) {
	return nil, errors.New("not found order")
}

func (o *orderRepository) UpdateStatus(orderID, newStatus string) error {
	return nil
}
