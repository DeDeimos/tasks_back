package ds

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	jwt.StandardClaims          // все что точно необходимо по RFC
	User_ID            uint     `json:"user_uuid"`            // наши данные - uuid этого пользователя в базе данных
	Scopes             []string `json:"scopes" json:"scopes"` // список доступов в нашей системе
	Role               string
}

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	ExpiresIn   time.Duration `json:"expires_in"`
	AccessToken string        `json:"access_token"`
	TokenType   string        `json:"token_type"`
	Role        string        `json:"role"`
}
