package auth

import (
	"{{.ModulePath}}/infra/jwt"
	"{{.ModulePath}}/logic/user"
)

type Service struct {
	userService *user.Service
}

func NewService() *Service {
	return &Service{
		userService: user.NewService(),
	}
}

func (s *Service) Login(username, password string) (*user.LoginResp, error) {
	req := &user.LoginReq{
		Username: username,
		Password: password,
	}

	loginUser, err := s.userService.Login(req)
	if err != nil {
		return nil, err
	}

	// 生成token
	accessToken, err := jwt.GenerateAccessToken(loginUser.ID, loginUser.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(loginUser.ID, loginUser.Username)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*user.LoginResp, error) {
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	_, err = s.userService.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// 生成新的token对
	accessToken, err := jwt.GenerateAccessToken(claims.UserID, claims.Username)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := jwt.GenerateRefreshToken(claims.UserID, claims.Username)
	if err != nil {
		return nil, err
	}

	return &user.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}