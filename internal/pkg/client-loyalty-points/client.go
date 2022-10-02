package client_loyalty_points

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

type ClientLoyaltyPoints struct {
	log        *zap.Logger
	address    string
	httpClient *http.Client
}

func NewClientLoyaltyPoints(log *zap.Logger, address string) *ClientLoyaltyPoints {
	return &ClientLoyaltyPoints{
		log:        log,
		address:    address,
		httpClient: &http.Client{},
	}
}

var ClientResponseOrderStatusRegistered = "REGISTERED"
var ClientResponseOrderStatusInvalid = "INVALID"
var ClientResponseOrderStatusProcessing = "PROCESSING"
var ClientResponseOrderStatusProcessed = "PROCESSED"

type ClientResponseOrderDTO struct {
	Order   string      `json:"order"`
	Status  string      `json:"status"`
	Accrual json.Number `json:"accrual,omitempty"`
}

func (c *ClientLoyaltyPoints) GetOrder(ctx context.Context, orderNumber string) (*ClientResponseOrderDTO, error) {
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		strings.Join([]string{c.address, "/api/orders/", orderNumber}, ""),
		nil,
	)

	if err != nil {
		return nil, err
	}

	response, err := c.httpClient.Do(request)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	c.log.Info("response status", zap.String("httpStatus", response.Status))

	responseBodyBytes, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	c.log.Info("response body", zap.ByteString("body", responseBodyBytes))

	var responseOrderDTO *ClientResponseOrderDTO

	if err := json.Unmarshal(responseBodyBytes, &responseOrderDTO); err != nil {
		return nil, err
	}

	return responseOrderDTO, nil
}
