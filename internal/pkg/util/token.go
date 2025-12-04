package util

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	"vibe-music-server/internal/config"
)

type Claims struct {
	Role     string `json:"role"`
	UserId   uint64 `json:"UserId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var (
	secret     string
	expiration int64
)

func init() {
	Jwt := config.Get().Jwt
	secret = Jwt.Secret
	expiration = Jwt.Expiration
}

func GenerateToken(claims Claims) string {
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Duration(expiration) * time.Minute))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
