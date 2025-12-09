package user

import (
	"strconv"
	"{{.ModulePath}}/logic/user"
	"{{.ModulePath}}/logic/shared"
	"{{.ModulePath}}/web/base"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService *user.Service
}

func NewHandler(userService *user.Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

// GetProfile 获取用户信息
func (h *Handler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		base.Error(c, 1000, "user not found in context")
		return
	}

	u, err := h.userService.GetByID(userID.(uint))
	if err != nil {
		if err == shared.ErrUserNotFound {
			base.Error(c, 2001, "user not found")
			return
		}
		base.Error(c, 500, err.Error())
		return
	}

	base.Success(c, u)
}

// GetUser 根据ID获取用户
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		base.Error(c, 400, "invalid user id")
		return
	}

	u, err := h.userService.GetByID(uint(id))
	if err != nil {
		if err == shared.ErrUserNotFound {
			base.Error(c, 2001, "user not found")
			return
		}
		base.Error(c, 500, err.Error())
		return
	}

	base.Success(c, u)
}