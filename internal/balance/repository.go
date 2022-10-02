package balance

import (
	"database/sql"
)

type balanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *balanceRepository {
	return &balanceRepository{db: db}
}

func (b *balanceRepository) Create(userID string) error {
	return nil
}

func (b *balanceRepository) FindByUser(userID string) (*BalanceEntity, error) {
	return &BalanceEntity{}, nil
}

func (b *balanceRepository) Accrue(userID string, amount int) error {
	return nil
}

func (b *balanceRepository) Withdraw(userID string, amount int) error {
	return nil
}
