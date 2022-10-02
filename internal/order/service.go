package order

type orderService struct {
	rep *orderRepository
}

func NewOrderService(rep *orderRepository) *orderService {
	return &orderService{
		rep: rep,
	}
}

func (o *orderService) GetOrdersByUser(userID string) ([]*Order, error) {
	orderEntities, err := o.rep.GetOrdersByUser(userID)

	if err != nil {
		return nil, err
	}

	orders := make([]*Order, 0, len(orderEntities))

	for i, orderEntity := range orderEntities {
		orders[i] = NewOrderFromEntity(orderEntity)
	}

	return orders, nil
}

func (o *orderService) AddOrder(userID, orderNumber string) (*OrderEntity, error) {
	orderEntity := OrderEntity{
		Number: orderNumber,
		Status: "REGISTERED",
		UserID: userID,
	}

	err := o.rep.Save(&orderEntity)

	if err != nil {
		return nil, err
	}

	return &orderEntity, nil
}

func (o *orderService) setStatus(orderID, newStatus string) error {
	//orderEntity, err := o.rep.FindByID(orderID)
	_, err := o.rep.FindByID(orderID)

	if err != nil {
		return err
	}

	// TODO: реализация конечного автомата для смены статуса
	// Или подумать и перенести его как метод в сущность OrderEntity.
	// Тогда мы сможем это сделать прямо в repository и в транзакционном формате

	if err := o.rep.UpdateStatus(orderID, newStatus); err != nil {
		return err
	}

	return nil
}

func (o *orderService) SetProcessStatusById(orderID string) error {
	err := o.setStatus(orderID, "PROCESSING")

	return err
}
