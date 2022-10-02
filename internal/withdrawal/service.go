package withdrawal

import (
	"context"
	"errors"
	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/pkg/luhn"
	"go.uber.org/zap"
	"strconv"
)

type WithdrawalService struct {
	log            *zap.Logger
	rep            *WithdrawalRepository
	balanceService *balance.BalanceService
}

func NewWithdrawalService(log *zap.Logger, rep *WithdrawalRepository, balanceService *balance.BalanceService) *WithdrawalService {
	return &WithdrawalService{
		log:            log,
		rep:            rep,
		balanceService: balanceService,
	}
}

// TODO: Нужва валидация параметров
func (s *WithdrawalService) Create(ctx context.Context, userID string, order string, sum int) error {
	orderNumber, err := strconv.Atoi(order)

	if err != nil {
		return ErrWithdrawalOrderNumberIsInvalid
	}

	valid := luhn.Valid(orderNumber)

	if !valid {
		return ErrWithdrawalOrderNumberIsInvalid
	}

	if err := s.rep.Create(ctx, userID, order, sum); err != nil {
		return err
	}

	errWithdraw := s.balanceService.Withdraw(ctx, userID, sum)

	if errWithdraw != nil {
		if errDeleteOrder := s.rep.deleteByOrder(ctx, order); errDeleteOrder != nil {
			return err
		}

		if errors.Is(errWithdraw, balance.ErrBalanceCannotBeNegative) {
			return ErrWithdrawalNotFoundsInBalance
		}

		return errWithdraw
	}

	return nil
}

func (s *WithdrawalService) GetAllByUser(ctx context.Context, userID string) ([]*WithdrawalEntity, error) {
	withdrawals, err := s.rep.FindAllByUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
