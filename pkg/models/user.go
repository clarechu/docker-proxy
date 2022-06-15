package models

type User struct {
	Account    string `json:"account,omitempty"`
	Password   string `json:"password,omitempty"`
	BasicToken string `json:"basicToken,omitempty"`
}
