package models

var (
	AuthorizationKey = "Authorization"

	WwwAuthorizationKey = "WWW-Authenticate"

	HttpSchema Schema = "http"

	HttpsSchema Schema = "https"
)

// Schema 协议
type Schema string

func (s Schema) SchemaToString() string {
	switch s {
	case HttpSchema:
		return "http"
	case HttpsSchema:
		return "https"
	}
	return ""
}
