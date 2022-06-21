package proxy

import (
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/clarechu/docker-proxy/pkg/utils/queue"
)

type App struct {
	Host                    string
	Token                   string
	Schema                  string
	LoggingHandler          models.LoggingHandler
	Stop                    chan struct{}
	Queue                   queue.Instance
	OAuth2EventHandlerFuncs models.OAuth2EventHandlerFuncs
}
