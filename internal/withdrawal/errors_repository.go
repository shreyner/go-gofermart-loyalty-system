package withdrawal

import "errors"

var ErrWithdrawalNotFound = errors.New("withdrawal not found")
var ErrWithdrawalOrderIsExist = errors.New("withdrawal with order is exist")
