package security

import (
	"github.com/golang-jwt/jwt/v4"
	"os/user"
	"time"
)

func issure() string {
	u, err := user.Current()
	if err != nil {
		return "jwt.issure"
	}
	return u.Username
}

var defaultPayload = map[string]interface{}{
	"exp": time.Now().Add(time.Minute * 20).Unix(),
	"iat": time.Now().Unix(),
	"nbf": time.Now().Unix(),
	"iss": issure(),
}

func GenerateJWTToken(
	tokenPayload, tokenSettings map[string]interface{},
	secretKey string, singMethod jwt.SigningMethod,
) (string, error) {
	claims := jwt.MapClaims{}

	for k, v := range defaultPayload {
		claims[k] = v
	}
	for k, v := range tokenSettings {
		claims[k] = v
	}
	for k, v := range tokenPayload {
		claims[k] = v
	}

	var sign jwt.SigningMethod
	if singMethod == nil {
		sign = jwt.SigningMethodHS256
	} else {
		sign = singMethod
	}

	token := jwt.NewWithClaims(sign, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
