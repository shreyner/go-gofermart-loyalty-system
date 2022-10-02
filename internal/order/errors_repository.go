package order

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderIsExist = errors.New("order is exist")
