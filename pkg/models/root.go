package models

type Root struct {
	Port     int32  `json:"port"`
	CertFile string `json:"certFile"`
	KeyFile  string `json:"keyFile"`
	Timeout  int32  `json:"timeout"`
	App      *App   `json:"app"`
}
