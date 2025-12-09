package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"{{.ModulePath}}/infra/dbs"
	"{{.ModulePath}}/logic/shared"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type RegisterReq struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *Service) Register(req *RegisterReq) error {
	// 检查用户是否已存在
	var count int64
	dbs.DB.Model(&User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		return shared.ErrUserExists
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Status:   shared.UserStatusActive,
	}

	return dbs.DB.Create(user).Error
}

func (s *Service) Login(req *LoginReq) (*User, error) {
	var user User
	err := dbs.DB.Where("username = ? AND status = ?", req.Username, shared.UserStatusActive).First(&user).Error
	if err != nil {
		return nil, shared.ErrUserNotFound
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, shared.ErrInvalidPassword
	}

	return &user, nil
}

func (s *Service) GetByID(id uint) (*User, error) {
	var user User
	err := dbs.DB.Where("id = ? AND status = ?", id, shared.UserStatusActive).First(&user).Error
	if err != nil {
		return nil, shared.ErrUserNotFound
	}
	return &user, nil
}