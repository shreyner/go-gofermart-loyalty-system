package user

import (
	"context"
)

type UserService struct {
	rep *userRepository
}

func NewUserService(rep *userRepository) *UserService {
	return &UserService{
		rep: rep,
	}
}

func (u *UserService) CreateUser(ctx context.Context, login, password string) (*User, error) {
	userEntity := UserEntity{Login: login}

	if err := userEntity.SetPassword(password); err != nil {
		return nil, err
	}

	if err := u.rep.Create(ctx, &userEntity); err != nil {
		return nil, err
	}

	user := CreateUserFromEntity(&userEntity)

	return user, nil
}

func (u *UserService) FindAndVerifyPassword(ctx context.Context, login, password string) (*User, error) {
	userEntity, err := u.rep.FindByLogin(ctx, login)

	if err != nil {
		return nil, err
	}

	if valid := userEntity.VerifyPassword(password); !valid {
		return nil, ErrUserPasswordIncorrect
	}

	return CreateUserFromEntity(userEntity), nil
}
