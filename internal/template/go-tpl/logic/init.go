package logic

import (
	"github.com/google/wire"
	"{{.ModulePath}}/logic/auth"
	"{{.ModulePath}}/logic/user"
)

var (
	UserService  *user.Service
	AuthService  *auth.Service
)

func Init() {
	wire.Build(ProviderSet)
}