package withdrawal

import "time"

type WithdrawalEntity struct {
	ID          string
	OrderNumber string
	Sum         int
	CreatedAt   time.Time
	UserID      string
}
