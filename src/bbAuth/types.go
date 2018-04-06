package bbAuth

// Token the representation of the OAuth2 token
type Token struct {
	Token           string  `json:"token"`
	RefreshToken    string  `json:"refreshToken"`
	ExpireEpochTime float64 `json:"expiration"`
}

// AuthenticationHandler a handler to manage the authentication
type AuthenticationHandler interface {
	// Manage the login via username/password
	UsernamePasswordLogin(username string, password string) (token Token, err error)
	// Manage the token refresh
	RefreshToken(refreshToken string) (newToken string, err error)
	// Manage the logout process
	Logout(token string) (loggedOut bool, err error)
}
