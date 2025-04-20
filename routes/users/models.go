package users

import "ChaiwalaBackend/db"

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GeneratedJWTResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	TokenType    string `json:"tokenType"`
}

type LoginUserResponse struct {
	Token GeneratedJWTResponse `json:"token"`
	User  db.User              `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
}
