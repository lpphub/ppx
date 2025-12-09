package auth

import (
	"{{.ModulePath}}/logic/auth"
	"{{.ModulePath}}/logic/shared"
	"{{.ModulePath}}/web/base"
	"{{.ModulePath}}/web/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	authService *auth.Service
}

func NewHandler(authService *auth.Service) *Handler {
	return &Handler{
		authService: authService,
	}
}

// Register 用户注册
func (h *Handler) Register(c *gin.Context) {
	var req types.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		base.Error(c, 400, err.Error())
		return
	}

	err := h.authService.Register(&auth.RegisterReq{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		if err == shared.ErrUserExists {
			base.Error(c, 2002, "user already exists")
			return
		}
		base.Error(c, 500, err.Error())
		return
	}

	base.Success(c, nil)
}

// Login 用户登录
func (h *Handler) Login(c *gin.Context) {
	var req types.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		base.Error(c, 400, err.Error())
		return
	}

	resp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		if err == shared.ErrUserNotFound {
			base.Error(c, 2001, "user not found")
			return
		}
		if err == shared.ErrInvalidPassword {
			base.Error(c, 2004, "invalid password")
			return
		}
		base.Error(c, 500, err.Error())
		return
	}

	base.Success(c, resp)
}

// RefreshToken 刷新token
func (h *Handler) RefreshToken(c *gin.Context) {
	var req types.RefreshTokenReq
	if err := c.ShouldBindJSON(&req); err != nil {
		base.Error(c, 400, err.Error())
		return
	}

	resp, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		base.Error(c, 1001, "invalid token")
		return
	}

	base.Success(c, resp)
}