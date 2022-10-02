package withdrawal

import (
	"context"
	"database/sql"
)

type withdrawalRepository struct {
	db *sql.DB
}

func NewWithdrawalRepository(db *sql.DB) *withdrawalRepository {
	return &withdrawalRepository{
		db: db,
	}
}

func (wr *withdrawalRepository) InitSchema(ctx context.Context) error {
	_, err := wr.db.ExecContext(
		ctx,
		`
			create table if not exists withdrawals
			(
				id         uuid        default gen_random_uuid() not null constraint withdrawals_pk primary key,
				"order"    integer                               not null constraint withdrawals_order_uniqk unique,
				sum        integer                               not null,
				user_id    uuid                                  not null constraint withdrawals_user_id_fk references users (id),
				created_at timestamptz default current_timestamp not null
			);
		`,
	)

	return err
}
