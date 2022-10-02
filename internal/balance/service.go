package balance

type balanceService struct {
	rep *balanceRepository
}

func NewBalanceService(rep *balanceRepository) *balanceService {
	return &balanceService{
		rep: rep,
	}
}

func (b *balanceService) CreateByUserID(userID string) error {
	return b.rep.Create(userID)
}

func (b *balanceService) GetByUserID(userID string) (*Balance, error) {
	balanceEntity, err := b.rep.FindByUser(userID)

	if err != nil {
		return nil, err
	}

	balance := &Balance{
		Current:   balanceEntity.Current,
		Withdrawn: balanceEntity.Withdrawn,
	}

	return balance, nil
}
