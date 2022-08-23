package authS

import (
	"blog-api/internal/dao"
)

type Auth struct {
	Username string
	Password string
}

func (a *Auth) Check() error {
	return dao.CheckAuth(a.Username, a.Password)
}
