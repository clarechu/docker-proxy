package models

type Root struct {
	Port    int32 `json:"port"`
	Timeout int32 `json:"timeout"`
	App     *App  `json:"app"`
}
