package order

import (
	"context"
	"errors"
	"go-gofermart-loyalty-system/pkg/luhn"
	"strconv"
)

type OrderService struct {
	rep *OrderRepository
}

func NewOrderService(rep *OrderRepository) *OrderService {
	return &OrderService{
		rep: rep,
	}
}

func (o *OrderService) GetOrdersByUser(ctx context.Context, userID string) ([]*OrderEntity, error) {
	return o.rep.GetOrdersByUser(ctx, userID)
}

func (o *OrderService) GetOrderByNumber(ctx context.Context, number string) (*OrderEntity, error) {
	return o.rep.FindByNumber(ctx, number)
}

func (o *OrderService) AddOrder(ctx context.Context, userID, orderNumber string) (*OrderEntity, error) {
	orderInt, err := strconv.Atoi(orderNumber)

	if err != nil {
		return nil, ErrOrderNumberIsInvalid
	}

	if !luhn.Valid(orderInt) {
		return nil, ErrOrderNumberIsInvalid
	}

	orderEntity, err := o.rep.Create(ctx, userID, orderNumber)

	if errors.Is(err, ErrOrderIsExist) {
		orderEntity, err := o.rep.FindByNumber(ctx, orderNumber)

		if err != nil {
			return nil, err
		}

		if orderEntity.UserID != userID {
			// TODO: Подумать как можно по другому обыграть ошибку с уже добавленным от другого пользователя
			return nil, ErrOrderAlreadyExistAnotherUser
		}

		return nil, ErrOrderIsExist
	}

	if err != nil {
		return nil, err
	}

	return orderEntity, nil
}

func (o *OrderService) setStatus(ctx context.Context, orderNumber, newStatus string) error {
	//orderEntity, err := o.rep.FindByID(orderNumber)
	//_, err := o.rep.FindByID(orderNumber)

	//if err != nil {
	//	return err
	//}

	// TODO: реализация конечного автомата для смены статуса
	// Или подумать и перенести его как метод в сущность OrderEntity.
	// Тогда мы сможем это сделать прямо в repository и в транзакционном формате

	if err := o.rep.UpdateStatusByOrderNumber(ctx, orderNumber, newStatus); err != nil {
		return err
	}

	return nil
}

func (o *OrderService) setStatusAndAccrual(ctx context.Context, orderNumber string, accrual float64) error {
	if err := o.rep.UpdateStatusAndAccuralByOrderNumber(ctx, orderNumber, StatusOrderProcessed, accrual); err != nil {
		return err
	}

	return nil
}

func (o *OrderService) SetProcessingStatusByNumber(ctx context.Context, orderNumber string) error {
	err := o.setStatus(ctx, orderNumber, StatusOrderProcessing)

	return err
}

func (o *OrderService) SetInvalidStatusByNumber(ctx context.Context, orderNumber string) error {
	err := o.setStatus(ctx, orderNumber, StatusOrderInvalid)

	return err
}

func (o *OrderService) SetProcessedStatusByNumber(ctx context.Context, orderNumber string, accrual float64) error {
	err := o.setStatusAndAccrual(ctx, orderNumber, accrual)

	return err
}
