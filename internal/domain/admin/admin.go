package admin

import "context"

type Admin struct {
	id       int32
	login    string
	password string
}

func (a *Admin) Password() string {
	return a.password
}

type AdminRepository interface {
	GetByLogin(ctx context.Context, login string) (*Admin, error)
}

func New(id int32, login, password string) *Admin {
	return &Admin{
		id:       id,
		login:    login,
		password: password,
	}
}
