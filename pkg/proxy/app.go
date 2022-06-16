package proxy

import "github.com/clarechu/docker-proxy/pkg/models"

type App struct {
	Host                    string
	Token                   string
	Schema                  string
	OAuth2EventHandlerFuncs models.OAuth2EventHandlerFuncs
}
