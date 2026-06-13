package admin

import "context"

type Admin struct {
	id           int32
	login        string
	passwordHash string
}

func (a *Admin) PasswordHash() string {
	return a.passwordHash
}

type AdminRepository interface {
	GetByLogin(ctx context.Context, login string) (*Admin, error)
}

func New(id int32, login, passwordHash string) *Admin {
	return &Admin{
		id:           id,
		login:        login,
		passwordHash: passwordHash,
	}
}
