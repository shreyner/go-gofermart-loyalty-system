package user

import (
	"database/sql"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(user UserEntity) error {
	return nil
}

func (u *userRepository) FindByLogin(login string) *UserEntity {
	return &UserEntity{ID: "1", Login: login}
}
