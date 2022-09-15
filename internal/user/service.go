package user

import "errors"

type userService struct {
	rep *userRepository
}

func NewUserService(rep *userRepository) *userService {
	return &userService{
		rep: rep,
	}
}

func (u *userService) CreateUser(login, password string) (*User, error) {
	userEntity := UserEntity{Login: login}

	if err := userEntity.SetPassword(password); err != nil {
		return nil, err
	}

	if err := u.rep.Create(userEntity); err != nil {
		return nil, err
	}

	user := CreateUserFromEntity(&userEntity)

	return user, nil
}

func (u *userService) FindAndVerifyPassword(login, password string) (*User, error) {
	userEntity := u.rep.FindByLogin(login)

	if userEntity == nil {
		return nil, errors.New("user not found")
	}

	if valid := userEntity.VerifyPassword(password); valid == false {
		return nil, errors.New("password incorrect")
	}

	return CreateUserFromEntity(userEntity), nil
}
