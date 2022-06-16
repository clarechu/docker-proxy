# docker-proxy

docker 代理设置

```go

import "github.com/clarechu/docker-proxy/pkg/models"
import "github.com/clarechu/docker-proxy/pkg/router"


func main() {
    root := &models.Root{
        App: &models.App{
        DockerRegistryHost: "localhost:9001",
        RegistryType:       models.NexusRegistry,
        Harbor: models.Harbor{
            Username: "admin",
            Password: "admin123",
        },
        Nexus: models.Nexus{
            Username: "admin",
            Password: "admin123",
        },
        Schema: models.HttpSchema,
        OAuth2EventHandlerFuncs: models.OAuth2EventHandlerFuncs{
                LoginFunc:      LoginFunc,
                CheckTokenFunc: CheckTokenFunc,
                PostTokenFunc:  PostTokenFunc,
            },
        },
        Port: 7777,    
    }   
	server := router.NewServer(root)
    go server.Run()
	
}



```

1. 重写LoginFunc 登陆handler

```go

// LoginFunc 登陆的function
func LoginFunc(user *models.User) (*models.OAuth2, error) {
	password := fmt.Sprintf("Basic %s", base64.EncodeToString("qwer:admin"))
	if user.Account == "qwer" && user.BasicToken == password {
		return &models.OAuth2{
			RefreshToken: "kas9Da81Dfa8",
			AccessToken:  "eyJhbGciOiJFUzI1NiIsInR5",
			ExpiresIn:    900,
		}, nil
	}
	return nil, errors.New("no token")
}
```

2. 重写CheckTokenFunc 验证token是否合法

```go
// CheckTokenFunc 校验token 是否合法
func CheckTokenFunc(token string) bool {
	if "Bearer eyJhbGciOiJFUzI1NiIsInR5" == token {
		return true
	}
	return false
}
```

3. 重写PostTokenFunc 更新和获取新的token

```go


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

```


* 注意当前仅支持http 协议暂时不支持http 的schema