package models

type App struct {
	DockerRegistryHost string `json:"dockerRegistryHost" default:"auth.docker.io"`
}
