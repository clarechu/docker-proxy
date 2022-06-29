package example

import (
	"github.com/clarechu/docker-proxy/pkg/models"
)

func NewNexusApp() *models.NexusApp {
	return &models.NexusApp{
		URL:      "http://localhost:8081",
		Username: "admin",
		Password: "admin123",
		Ports: []models.Port{
			{
				Name:       "docker-hosted",
				Port:       9000,
				TargetPort: 9001,
			},
		},
		Schema:         models.HttpSchema,
		LoggingHandler: LoggingHandler,
		OAuth2EventHandlerFuncs: models.OAuth2EventHandlerFuncs{
			LoginFunc:      LoginFunc,
			CheckTokenFunc: CheckTokenFunc,
			PostTokenFunc:  PostTokenFunc,
		},
	}
}
