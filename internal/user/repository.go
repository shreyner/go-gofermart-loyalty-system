package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(ctx context.Context, user *UserEntity) error {
	row := u.db.QueryRowContext(
		ctx,
		`insert into users (login, password) values ($1, $2) returning id;`,
		user.Login,
		user.password,
	)

	if row.Err() != nil {
		var pgErr *pgconn.PgError

		if !errors.As(row.Err(), &pgErr) {
			return row.Err()
		}

		if pgErr.Code == pgerrcode.UniqueViolation && pgErr.ConstraintName == "users_login_key" {
			return fmt.Errorf("%q: %w", user.Login, ErrLoginAlreadyExist)
		}

		return row.Err()
	}

	if err := row.Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (u *userRepository) FindByLogin(ctx context.Context, login string) (*UserEntity, error) {
	row := u.db.QueryRowContext(
		ctx,
		`select id, login, password from users where login = $1 limit 1;`,
		login,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	user := UserEntity{}

	err := row.Scan(&user.ID, &user.Login, &user.password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
