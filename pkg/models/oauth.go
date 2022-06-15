package models

type OAuth2 struct {
	// RefreshToken (Optional) Token which can be used to get additional access tokens for the same subject with different scopes.
	// This token should be kept secure by the client and only sent to the authorization server which issues bearer tokens.
	// This field will only be set when `access_type=offline` is provided in the request.
	RefreshToken string `json:"refresh_token"`
	// IssuedAt (Optional) The RFC3339-serialized UTC standard time at which a given token was issued.
	// If issued_at is omitted, the expiration is from when the token exchange completed.
	IssuedAt string `json:"issued_at"`
	// ExpiresIn (REQUIRED) The duration in seconds since the token was issued that it will remain valid.
	// When omitted, this defaults to 60 seconds.
	// For compatibility with older clients, a token should never be returned with less than 60 seconds to live.
	ExpiresIn string `json:"expires_in"`
	// Scope (REQUIRED) The scope granted inside the access token.
	// This may be the same scope as requested or a subset. This requirement is stronger than specified in
	// [RFC6749 Section 4.2.2](https://tools.ietf.org/html/rfc6749#section-4.2.2)
	// by strictly requiring the scope in the return value.
	Scope string `json:"scope"`
	// AccessToken (REQUIRED) An opaque Bearer token that clients should supply to subsequent requests in the Authorization header.
	//This token should not be attempted to be parsed or understood by the client but treated as opaque string.
	AccessToken string `json:"access_token"`
}

type Token struct {
	// GrantType (REQUIRED) Type of grant used to get token.
	// When getting a refresh token using credentials this type should be set to "password" and
	// have the accompanying username and password parameters.
	// Type "authorization_code" is reserved for future use for authenticating to an authorization server without having to send credentials directly from the client.
	// When requesting an access token with a refresh token this should be set to "refresh_token".
	// password 、 authorization_code 、 refresh_token
	GrantType string `json:"grant_type" schema:"grant_type,required"`
	// Service (REQUIRED) The name of the service which hosts the resource to get access for.
	// Refresh tokens will only be good for getting tokens for this service.
	Service string `json:"service" schema:"service,required"`
	// ClientID (REQUIRED) String identifying the client.
	// This client_id does not need to be registered with the authorization server but should be set to
	// a meaningful value in order to allow auditing keys created by unregistered clients.
	// Accepted syntax is defined in [RFC6749 Appendix A.1](https://tools.ietf.org/html/rfc6749#appendix-A.1)
	ClientID string `json:"client_id" schema:"client_id,required"`
	// AccessType (OPTIONAL) Access which is being requested. If "offline" is provided then a refresh token will be returned.
	// The default is "online" only returning short lived access token.
	// If the grant type is "refresh_token" this will only return the same refresh token and not a new one.
	AccessType string `json:"access_type" schema:"access_type"`
	// Scope (OPTIONAL) The resource in question,
	// formatted as one of the space-delimited entries from the scope parameters from the WWW-Authenticate header shown above.
	// This query parameter should only be specified once but may contain multiple scopes
	// using the scope list format defined in the scope grammar. If multiple scope is provided from
	// WWW-Authenticate header the scopes should first be converted to a scope list before requesting the token.
	// The above example would be specified as: scope=repository:samalba/my-app:push.
	// When requesting a refresh token the scopes may be empty since the refresh token will not be limited by this scope,
	// only the provided short lived access token will have the scope limitation.
	Scope string `json:"scope" schema:"scope"`
	// RefreshToken (OPTIONAL) The refresh token to use for authentication when grant type "refresh_token" is used.
	RefreshToken string `json:"refresh_token" schema:"refresh_token"`
	// Username (OPTIONAL) The username to use for authentication when grant type "password" is used.
	Username string `json:"username" schema:"username"`
	// Password (OPTIONAL) The password to use for authentication when grant type "password" is used.
	Password string `json:"password" schema:"password"`
}

/*
client_id= docker

grant_type = refresh_token

refresh_token = kas9Da81Dfa8

scope = repository%3Asonatype%2Fnexus3%3Apush%2Cpull

service = http%3A%2F%2F172.20.218.78%3A7777%2Fv2%2Ftoken
*/
