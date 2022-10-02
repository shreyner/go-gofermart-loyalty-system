package order

import (
	"context"
	"database/sql"
	"errors"
)

type orderRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *orderRepository {
	return &orderRepository{db: db}
}

func (o *orderRepository) InitSchema(ctx context.Context) error {
	_, err := o.db.ExecContext(
		ctx,
		`
			create table if not exists orders
			(
				id         uuid        default gen_random_uuid() not null constraint orders_pk unique primary key,
				number     integer                               not null constraint orders_number_uniqk unique,
				status     varchar,
				accrual    integer,
				user_id    uuid                                  not null constraint orders_users_null_fk references users (id),
				created_at timestamptz default current_timestamp
			);
		`,
	)

	return err
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
