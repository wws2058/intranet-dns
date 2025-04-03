package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tswcbyy1107/intranet-dns/config"
	"github.com/tswcbyy1107/intranet-dns/utils"
)

type AppJwtClaim struct {
	Username string `json:"username,omitempty"`
	jwt.StandardClaims
}

const tokenExpireDuration = 24 * time.Hour

// get jwt token
func GenJwtToken(username string) (string, error) {
	c := AppJwtClaim{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
			Issuer:    config.GlobalConfig.App.Name,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(utils.JwtSecret)
}

// parse jwt token
func ParseToken(tokenString string) (*AppJwtClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AppJwtClaim{}, func(token *jwt.Token) (i interface{}, err error) {
		return utils.JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AppJwtClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
