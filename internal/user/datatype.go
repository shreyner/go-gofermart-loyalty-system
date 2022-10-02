package user

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
	return nil
}

func (u *UserEntity) VerifyPassword(password string) bool {
	return true
}
