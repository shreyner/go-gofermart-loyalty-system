package balance

import "context"

type BalanceService struct {
	rep *balanceRepository
}

func NewBalanceService(rep *balanceRepository) *BalanceService {
	return &BalanceService{
		rep: rep,
	}
}

func (b *BalanceService) CreateByUserID(ctx context.Context, userID string) error {
	return b.rep.Create(ctx, userID)
}

func (b *BalanceService) GetByUserID(ctx context.Context, userID string) (*Balance, error) {
	balanceEntity, err := b.rep.FindByUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	balance := &Balance{
		Current:   balanceEntity.Current,
		Withdrawn: balanceEntity.Withdrawn,
	}

	return balance, nil
}

func (b *BalanceService) Accrue(ctx context.Context, userID string, amount int) error {
	return b.rep.Accrue(ctx, userID, amount)
}

func (b *BalanceService) Withdraw(ctx context.Context, userID string, amount int) error {
	return b.rep.Withdraw(ctx, userID, amount)
}
