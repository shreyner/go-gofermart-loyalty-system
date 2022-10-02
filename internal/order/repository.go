package order

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type OrderRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) InitSchema(ctx context.Context) error {
	_, err := o.db.ExecContext(
		ctx,
		`
			create table if not exists orders
			(
				id         uuid        default gen_random_uuid() not null constraint orders_pk unique primary key,
				number     varchar                               not null constraint orders_number_uniqk unique,
				status     varchar,
				accrual    integer,
				user_id    uuid                                  not null constraint orders_users_null_fk references users (id),
				created_at timestamptz default current_timestamp
			);
		`,
	)

	return err
}

func (o *OrderRepository) GetOrdersByUser(ctx context.Context, userID string) ([]*OrderEntity, error) {
	rows, err := o.db.QueryContext(
		ctx,
		`
			select id, user_id, number, status, accrual, created_at
			from orders
			where user_id = $1
			order by created_at;
		`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orders := make([]*OrderEntity, 0)

	for rows.Next() {
		var order OrderEntity

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Number,
			&order.Status,
			&order.Accrual,
			&order.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return orders, nil
}

func (o *OrderRepository) Create(ctx context.Context, userID string, number string) (*OrderEntity, error) {
	row := o.db.QueryRowContext(
		ctx,
		`
			insert into orders (user_id, number, status)
			values ($1, $2, $3)
			returning id, number, status, user_id, created_at;
		`,
		userID,
		number,
		StatusOrderNew,
	)

	if row.Err() != nil {
		var pgErr *pgconn.PgError

		if errors.As(row.Err(), &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, ErrOrderIsExist
		}

		return nil, row.Err()
	}

	order := OrderEntity{}

	err := row.Scan(
		&order.ID,
		&order.Number,
		&order.Status,
		&order.UserID,
		&order.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *OrderRepository) FindByNumber(ctx context.Context, number string) (*OrderEntity, error) {
	// TODO: Что делать с такими бесконечно длинными пречислениями в запроса? Или делать выборку строго по каким-то параметрам
	row := o.db.QueryRowContext(
		ctx,
		`select id, user_id, number, status, accrual, created_at from orders where number = $1`,
		number,
	)

	if errors.Is(row.Err(), sql.ErrNoRows) {
		return nil, ErrOrderNotFound
	}

	if row.Err() != nil {
		return nil, row.Err()
	}

	var orderEntity OrderEntity

	err := row.Scan(
		&orderEntity.ID,
		&orderEntity.UserID,
		&orderEntity.Number,
		&orderEntity.Status,
		&orderEntity.Accrual,
		&orderEntity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &orderEntity, err
}

func (o *OrderRepository) Save(orderEntity *OrderEntity) error {
	return nil
}

func (o *OrderRepository) FindByID(orderID string) (*OrderEntity, error) {
	return nil, errors.New("not found order")
}

func (o *OrderRepository) UpdateStatusByOrderNumber(ctx context.Context, number, newStatus string) error {
	_, err := o.db.ExecContext(
		ctx,
		`update orders set status = $1 where number = $2`,
		newStatus,
		number,
	)

	if err != nil {
		return err
	}

	return nil
}

// TODO: need refactoring
func (o *OrderRepository) UpdateStatusAndAccuralByOrderNumber(
	ctx context.Context,
	number,
	newStatus string,
	accural float64,
) error {
	_, err := o.db.ExecContext(
		ctx,
		`update orders set status = $1, accrual = $3 where number = $2`,
		newStatus,
		number,
		int(accural*100), // TODO: Hot fix and need refactoring
	)

	if err != nil {
		return err
	}

	return nil
}
