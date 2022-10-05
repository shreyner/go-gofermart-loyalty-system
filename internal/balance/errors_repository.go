package balance

import "errors"

var ErrUserHasBalance = errors.New("user has balance")
var ErrBalanceNotFound = errors.New("balance not found")
var ErrBalanceCannotBeNegative = errors.New("balance cannot be negative")
