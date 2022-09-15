package client_loyalty_points

import "errors"

type ClientLoyaltyPoints struct {
}

func (c *ClientLoyaltyPoints) GetOrder(orderNumber string) error {
	return errors.New("clientLoyaltyPoints: method not implements")
}
