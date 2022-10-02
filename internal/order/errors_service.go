package order

import "errors"

var ErrOrderNumberIsInvalid = errors.New("order number is invalid")
var ErrOrderAlreadyExistAnotherUser = errors.New("order is already exists another user")
