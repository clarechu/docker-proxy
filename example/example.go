package example

import (
	"errors"
	"fmt"
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/clarechu/docker-proxy/pkg/utils/base64"
)

func NewApp1() *models.App {
	return &models.App{
		RegistryType: models.NexusRegistry,
		Harbor: models.Harbor{
			Username: "admin",
			Password: "admin123",
		},
		Schema: models.HttpSchema,
		OAuth2EventHandlerFuncs: models.OAuth2EventHandlerFuncs{
			LoginFunc:      LoginFunc,
			CheckTokenFunc: CheckTokenFunc,
			PostTokenFunc:  PostTokenFunc,
		},
	}
}

// LoginFunc 登陆的function
func LoginFunc(user *models.User) (*models.OAuth2, error) {
	password := fmt.Sprintf("Basic %s", base64.EncodeToString("admin:admin"))
	if user.Account == "admin" && user.BasicToken == password {
		return &models.OAuth2{
			RefreshToken: "kas9Da81Dfa8",
			AccessToken:  "eyJhbGciOiJFUzI1NiIsInR5",
			ExpiresIn:    900000,
		}, nil
	}
	return nil, errors.New("no token")
}

// CheckTokenFunc 校验token 是否合法
func CheckTokenFunc(token string) bool {
	if "Bearer eyJhbGciOiJFUzI1NiIsInR5" == token {
		return true
	}
	return false
}

// PostTokenFunc 获取token令牌 和 刷新令牌
func PostTokenFunc(token *models.Token) (*models.OAuth2, error) {
	if token.RefreshToken == "kas9Da81Dfa8" {
		return &models.OAuth2{
			RefreshToken: "kas9Da81Dfa8",
			AccessToken:  "eyJhbGciOiJFUzI1NiIsInR5",
			ExpiresIn:    900000,
		}, nil
	}
	return nil, errors.New("get oauth error ")
}
