package auth

import (
	userType "go-gofermart-loyalty-system/internal/user"
)

type userService interface {
	CreateUser(login, password string) (*userType.User, error)
	FindAndVerifyPassword(login, password string) (*userType.User, error)
}

type balanceService interface {
	CreateByUserID(userID string) error
}

type authService struct {
	userService    userService
	balanceService balanceService
}

func NewAuthService(userService userService, balanceService balanceService) *authService {
	return &authService{
		userService:    userService,
		balanceService: balanceService,
	}
}

func (a *authService) RegistryByLogin(login, password string) (*userType.User, error) {
	user, err := a.userService.CreateUser(login, password)

	if err != nil {
		return nil, err
	}

	// TODO: Как тут решить вопрос с транзакционностью?
	// TODO: Подумать как ре организовать эту логку, так как при добавлении авторизации через Google, Yandex.ID будет дублирвоание логики
	if err := a.balanceService.CreateByUserID(user.ID); err != nil {
		return nil, err // TODO: Добавить свою ошибку
	}

	return user, nil
}

func (a *authService) Login(login, password string) (*userType.User, error) {
	user, err := a.userService.FindAndVerifyPassword(login, password)

	if err != nil {
		return nil, err
	}

	return user, err
}
