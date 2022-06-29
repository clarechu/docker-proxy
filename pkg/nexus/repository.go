package nexus

import (
	"github.com/clarechu/docker-proxy/pkg/models"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/client"
	"github.com/datadrivers/go-nexus-client/nexus3/pkg/repository"
	log "k8s.io/klog/v2"
)

type Repository struct {
	repositoryService *repository.RepositoryService
	client            *client.Client
}

func NewRepository(app models.NexusApp) *Repository {
	config := client.Config{
		Insecure: false,
		Password: app.Password,
		Username: app.Username,
		URL:      app.URL,
	}
	nexusClient := client.NewClient(config)
	return &Repository{
		client:            nexusClient,
		repositoryService: repository.NewRepositoryService(nexusClient),
	}
}

var (
	RepositoryFormat     = "docker"
	RepositoryGroupType  = "group"
	RepositoryHostedType = "hosted"
)

func (r *Repository) GetPortByDocker() ([]int, error) {
	ports := make([]int, 0)
	infos, err := r.repositoryService.List()
	if err != nil {
		return ports, err
	}
	for _, info := range infos {
		if info.Format == RepositoryFormat {
			switch info.Type {
			case RepositoryGroupType:
				groupService := r.repositoryService.Docker.Group
				group, err := groupService.Get(info.Name)
				if err != nil {
					return ports, err
				}
				port := group.Docker.HTTPPort
				ports = append(ports, *port)
				log.Infof("docker group name :%s, port %d", info.Name, port)
			case RepositoryHostedType:
				hostedService := r.repositoryService.Docker.Hosted
				hosted, err := hostedService.Get(info.Name)
				if err != nil {
					return ports, err
				}
				port := hosted.Docker.HTTPPort
				ports = append(ports, *port)

				log.Infof("docker hosted name :%s, port %d", info.Name, port)

			}
		}
	}
	return ports, err
}
