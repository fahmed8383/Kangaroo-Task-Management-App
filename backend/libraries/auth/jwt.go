package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims holds the claims for the jwt, standardclaims is added to create expiry date field
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJwt is responsible for creating a valid jwt for the logged in user
func GenerateJwt(username string, secret string) (string, error) {
	var claims Claims
	claims.Username = username
	claims.ExpiresAt = time.Now().Add(10 * time.Minute).Unix()
	unsignedtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := unsignedtoken.SignedString([]byte(secret))
	return token, err
}

// GetToken is responsible for making sure that the jwt has not been tampered and returns its values
// Does not check to see if token has expired
func GetToken(jwtString string, secret string) (*Claims, *jwt.Token, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(jwtString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	// get validation errors
	v, _ := err.(*jwt.ValidationError)

	// if there is any kind of error other than ValidationErrorExpired, return false and cancel request
	if err != nil && v.Errors != jwt.ValidationErrorExpired {
		return nil, nil, err
	}
	return claims, token, nil
}
