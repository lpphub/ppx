package logic

import (
	"github.com/google/wire"
	"{{.ModulePath}}/logic/auth"
	"{{.ModulePath}}/logic/user"
)

var ProviderSet = wire.NewSet(
	user.NewService,
	auth.NewService,
)