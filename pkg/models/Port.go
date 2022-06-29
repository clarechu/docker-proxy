package models

type Port struct {
	Name       string `json:"name,omitempty"`
	Port       int32  `json:"port,omitempty"`
	TargetPort int32  `json:"targetPort,omitempty"`
	Protocol   string `json:"protocol,omitempty"`
}
