package withdrawal

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
)

type WithdrawalRepository struct {
	log *zap.Logger
	db  *sql.DB
}

func NewWithdrawalRepository(log *zap.Logger, db *sql.DB) *WithdrawalRepository {
	return &WithdrawalRepository{
		log: log,
		db:  db,
	}
}

func (wr *WithdrawalRepository) Create(ctx context.Context, userID string, order string, sum int) error {
	_, err := wr.db.ExecContext(
		ctx,
		`insert into withdrawals ("order", sum, user_id) values ($1, $2 ,$3);`,
		order,
		sum,
		userID,
	)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
		return ErrWithdrawalOrderIsExist
	}

	if err != nil {
		return err
	}

	return nil
}

// TODO: Временно и нужно до того как решиться вопрос с транзакционностью при списании
func (wr *WithdrawalRepository) deleteByOrder(ctx context.Context, order string) error {
	result, err := wr.db.ExecContext(
		ctx,
		`delete from withdrawals where "order"=$1;`,
		order,
	)

	if err != nil {
		return err
	}

	rowAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowAffected == 0 {
		return ErrWithdrawalNotFound
	}

	return nil
}

func (wr *WithdrawalRepository) FindAllByUser(ctx context.Context, userID string) ([]*WithdrawalEntity, error) {
	rows, err := wr.db.QueryContext(
		ctx,
		`select id, user_id, "order", sum, created_at from withdrawals where user_id = $1 order by created_at;`,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	withdrawals := make([]*WithdrawalEntity, 0)

	for rows.Next() {
		var withdrawal WithdrawalEntity

		err := rows.Scan(
			&withdrawal.ID,
			&withdrawal.UserID,
			&withdrawal.OrderNumber,
			&withdrawal.Sum,
			&withdrawal.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, &withdrawal)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return withdrawals, nil
}
