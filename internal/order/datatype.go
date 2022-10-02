package order

type Order struct {
	Number     string
	Status     string
	Accrual    string
	UploadedAt string
}

func NewOrderFromEntity(orderEntity *OrderEntity) *Order {
	return &Order{
		Number:     orderEntity.Number,
		Status:     orderEntity.Status,
		Accrual:    orderEntity.Accrual,
		UploadedAt: orderEntity.CreatedAt,
	}
}

type OrderEntity struct {
	ID        string
	Number    string
	Status    string // TODO: Как типизировать enum
	Accrual   string
	CreatedAt string
	UserID    string
}
