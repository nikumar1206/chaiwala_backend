package users

type RegisterUser struct {
	Username string
	Password string
}

type LoginUser struct {
	Username string
	Password string
}

type GeneratedJWTResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
