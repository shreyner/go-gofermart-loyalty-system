package order

import "context"

type AsyncProcessingOrderService struct {
	service *OrderService
	pool    *WorkerPool
}

func NewAsyncProcessingOrderService(service *OrderService, pool *WorkerPool) *AsyncProcessingOrderService {
	return &AsyncProcessingOrderService{
		service: service,
		pool:    pool,
	}
}

func (a *AsyncProcessingOrderService) CreateOrderAndAddQueue(ctx context.Context, userID, orderNumber string) error {
	orderEntity, err := a.service.AddOrder(ctx, userID, orderNumber)

	if err != nil {
		return err
	}

	job := OrderJob{orderNumber: orderEntity.Number}
	a.pool.AddJob(&job)

	return nil
}
