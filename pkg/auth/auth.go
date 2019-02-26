package auth

const (
	// AuthHeader defines authorization header.
	AuthHeader = "Authorization"
	// AuthScheme defines authorization scheme.
	AuthScheme = "Bearer"
	// AuthorizationKey is the key used to store authorization token data
	AuthorizationKey = "authorization"
)

type AuthType int

const (
	NoAuth AuthType = iota
	JWTAuth
	APIKeyAuth
	HmacAuth
)
