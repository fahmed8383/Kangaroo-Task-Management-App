package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fahmed8383/SchedulingApp/libraries/api"
)

// Claims holds the claims for the jwt, standardclaims is added to create expiry date field
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJwt is responsible for creating a valid jwt for the logged in user
func GenerateJwt(user api.LoginInfo, secret string) (string, error) {
	var claims Claims
	claims.Username = user.Username
	claims.ExpiresAt = time.Now().Add(10 * time.Minute).Unix()
	unsignedtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := unsignedtoken.SignedString([]byte(secret))
	return token, err
}
