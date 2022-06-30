package models

type NexusApp struct {
	Host                    string                  `json:"host"`
	URL                     string                  `json:"url" default:"nexus.com"  validate:"required"`
	Ports                   []Port                  `json:"ports"`
	Username                string                  `json:"username,omitempty"`
	Password                string                  `json:"password,omitempty"`
	Schema                  Schema                  `json:"schema,omitempty"`
	OAuth2EventHandlerFuncs OAuth2EventHandlerFuncs `validate:"required"`
	LoggingHandler          LoggingHandler          `validate:"required"`
}

type App struct {
	RegistryType            RegistryType            `json:"registryType,omitempty" default:"harbor" validate:"required"`
	DockerRegistryHost      string                  `json:"dockerRegistryHost" default:"auth.docker.io"  validate:"required"`
	Nexus                   Nexus                   `json:"nexus,omitempty"`
	Harbor                  Harbor                  `json:"harbor,omitempty"`
	DockerRegistry          Docker                  `json:"dockerRegistry,omitempty"`
	Schema                  Schema                  `json:"schema,omitempty"`
	OAuth2EventHandlerFuncs OAuth2EventHandlerFuncs `validate:"required"`
	LoggingHandler          LoggingHandler          `validate:"required"`
}

type OAuth2EventHandlerFuncs struct {
	LoginFunc      `validate:"required"`
	CheckTokenFunc `validate:"required"`
	PostTokenFunc  `validate:"required"`
}

type LoggingHandler func(logging *Logging) error

// LoginFunc 登陆的function
type LoginFunc func(user *User) (*OAuth2, error)

// CheckTokenFunc 校验token 是否合法
type CheckTokenFunc func(token string) bool

// PostTokenFunc 获取token令牌 和 刷新令牌
type PostTokenFunc func(token *Token) (*OAuth2, error)

type Nexus struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Harbor struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Docker struct {
}

var (
	NexusRegistry  RegistryType = "nexus"
	HarborRegistry RegistryType = "harbor"
	DockerRegistry RegistryType = "docker"
)

type RegistryType string
