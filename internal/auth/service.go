package auth

import (
	"context"
	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/internal/user"
)

type AuthService struct {
	userService    *user.UserService
	balanceService *balance.BalanceService
}

func NewAuthService(userService *user.UserService, balanceService *balance.BalanceService) *AuthService {
	return &AuthService{
		userService:    userService,
		balanceService: balanceService,
	}
}

func (a *AuthService) RegistryByLogin(ctx context.Context, login, password string) (*user.User, error) {
	u, err := a.userService.CreateUser(ctx, login, password)

	if err != nil {
		return nil, err
	}

	// TODO: Как тут решить вопрос с транзакционностью?
	// TODO: Подумать как ре организовать эту логку, так как при добавлении авторизации через Google, Yandex.ID будет дублирвоание логики
	err = a.balanceService.CreateByUserID(ctx, u.ID)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a *AuthService) Login(ctx context.Context, login, password string) (*user.User, error) {
	u, err := a.userService.FindAndVerifyPassword(ctx, login, password)

	if err != nil {
		return nil, err
	}

	return u, err
}
