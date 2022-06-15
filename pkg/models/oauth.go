package models

type OAuth2 struct {
	RefreshToken string `json:"refresh_token"`
	IssuedAt     string `json:"issued_at"`
	ExpiresIn    string `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
}

type Token struct {
	// GrantType
	GrantType string `json:"grant_type" schema:"grant_type,required"`
	// Service
	Service string `json:"service" schema:"service,required"`
	// ClientID
	ClientID string `json:"client_id" schema:"client_id,required"`
	// AccessType
	AccessType string `json:"access_type" schema:"access_type"`
	// Scope
	Scope string `json:"scope" schema:"scope"`
	// RefreshToken
	RefreshToken string `json:"refresh_token" schema:"refresh_token"`
	// Username
	Username string `json:"username" schema:"username"`
	// Password
	Password string `json:"password" schema:"password"`
}

/*
client_id= docker

grant_type = refresh_token

refresh_token = kas9Da81Dfa8

scope = repository%3Asonatype%2Fnexus3%3Apush%2Cpull

service = http%3A%2F%2F172.20.218.78%3A7777%2Fv2%2Ftoken
*/
