package order

import (
	"database/sql"
	"time"
)

type OrderEntity struct {
	ID        string
	Number    string
	Status    string // TODO: Как типизировать enum
	Accrual   sql.NullInt64
	CreatedAt time.Time
	UserID    string
}
