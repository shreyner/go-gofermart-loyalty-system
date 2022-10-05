package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID    string `json:"id"`
	Login string `json:"login"`
}

func CreateUserFromEntity(user *UserEntity) *User {
	return &User{
		ID:    user.ID,
		Login: user.Login,
	}
}

type UserEntity struct {
	ID       string
	Login    string
	password string
}

func (u *UserEntity) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.password = string(hashedPassword)

	return nil
}

func (u *UserEntity) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false
	}

	if err != nil {
		return false
	}

	return true
}
