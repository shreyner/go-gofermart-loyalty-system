package order

type AsyncProcessingOrder struct {
	service *orderService
	pool    *WorkerPool
}

func NewAsyncProcessingOrder(service *orderService, pool *WorkerPool) *AsyncProcessingOrder {
	return &AsyncProcessingOrder{
		service: service,
		pool:    pool,
	}
}

func (a *AsyncProcessingOrder) CreateOrderAndAddQueue(userID, orderNumber string) error {

	orderEntity, err := a.service.AddOrder(userID, orderNumber)

	if err != nil {
		return err
	}

	job := OrderJob{orderNumber: orderEntity.Number}
	a.pool.AddJob(&job)

	return nil
}
