package balance

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type balanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *balanceRepository {
	return &balanceRepository{db: db}
}

func (b *balanceRepository) InitSchema(ctx context.Context) error {
	_, err := b.db.ExecContext(
		ctx,
		`
			create table if not exists balances
			(
				user_id    uuid constraint balances_pk unique constraint balances_users_null_fk references users (id),
				current    integer default 0,
				withdrawal integer default 0,
				created_at timestamptz default current_timestamp not null
			);
		`,
	)

	return err
}

func (b *balanceRepository) Create(ctx context.Context, userID string) error {
	_, err := b.db.ExecContext(
		ctx,
		`insert into balances (user_id) values ($1)`,
		userID,
	)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return ErrUserHasBalance
	}

	return err
}

func (b *balanceRepository) FindByUser(ctx context.Context, userID string) (*BalanceEntity, error) {
	row := b.db.QueryRowContext(
		ctx,
		`select user_id, "current", withdrawal from balances where user_id = $1`,
		userID,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	balanceEntity := BalanceEntity{}

	err := row.Scan(&balanceEntity.UserID, &balanceEntity.Current, &balanceEntity.Withdrawn)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBalanceNotFound
		}

		return nil, err
	}

	return &balanceEntity, nil
}

func (b *balanceRepository) Accrue(ctx context.Context, userID string, amount int) error {
	result, err := b.db.ExecContext(
		ctx,
		`update balances set current = current + $1 where user_id = $2`,
		amount,
		userID,
	)

	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrBalanceNotFound
	}

	return nil
}

func (b *balanceRepository) Withdraw(ctx context.Context, userID string, amount int) error {
	result, err := b.db.ExecContext(
		ctx,
		`
				update balances
				set current    = current - $1,
					withdrawal = withdrawal + $1
				where user_id = $2
				  and current - $1 >= 0;
		`,
		amount,
		userID,
	)

	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrBalanceCannotBeNegative
	}

	return nil
}
