package models

type App struct {
	RegistryType       RegistryType `json:"registryType,omitempty" default:"harbor"`
	DockerRegistryHost string       `json:"dockerRegistryHost" default:"auth.docker.io"`
	Nexus              Nexus        `json:"nexus,omitempty"`
	Harbor             Harbor       `json:"harbor,omitempty"`
	DockerRegistry     Docker       `json:"dockerRegistry,omitempty"`
	Schema             Schema       `json:"schema,omitempty"`
}

type LoginFunc func(user User) (OAuth2, error)

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
