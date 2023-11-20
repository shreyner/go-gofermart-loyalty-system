package withdrawal

import "errors"

var ErrWithdrawalOrderNumberIsInvalid = errors.New("order number is invalid")
var ErrWithdrawalNotFoundsInBalance = errors.New("not funds in balance")
