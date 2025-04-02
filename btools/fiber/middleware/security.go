package middleware

import (
	"errors"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/seemyown/backend-toolkit/btools/exc"
	"reflect"
	"strings"
)

var jwtLog = log.NewSubLogger("jwt_middleware")

type JwtMiddlewareConfig struct {
	Secret      string
	AuthKeyName string
	TokenType   string
	Issuer      string
	Out         interface{}
}

func JWTMiddleware(config *JwtMiddlewareConfig) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		headerValue := ctx.Get(config.AuthKeyName)
		if headerValue == "" {
			jwtLog.Warn("no auth header found")
			return ctx.Next()
		}

		var tokenString string
		if config.TokenType == "Bearer" {
			t, err := extractToken(headerValue)
			if err != nil {
				jwtLog.Error(err, "error extracting token")
				return exc.ForbiddenError("bad_token", err.Error())
			}
			tokenString = t
		} else {
			tokenString = headerValue
		}
		jwtLog.Info("Incoming request with token: %s", tokenString)
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(config.Secret), nil
		})

		if err != nil || !token.Valid {
			jwtLog.Error(err, "Incorrect token")
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				return exc.ForbiddenError("invalid_token", "wrong token signature")
			} else if "Token is expired" == err.Error() {
				return exc.UnauthorizedError("token_expired", "token is expired")
			}
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			jwtLog.Error(err, "Incorrect token")
			return exc.ForbiddenError("invalid_token", "bad claims")
		}

		if claims["iss"] != config.Issuer {
			jwtLog.Error(errors.New("wrong issuer"), "Incorrect issuer")
			return exc.ForbiddenError("wrong_issuer", "Wrong token issuer")
		}

		mappedData := reflect.New(reflect.TypeOf(config.Out)).Interface()
		if err := mapstructure.Decode(claims, mappedData); err != nil {
			jwtLog.Error(err, "Token mapping error")
			return exc.ForbiddenError("invalid_token", "claims mapping error")
		}

		v := reflect.ValueOf(mappedData).Elem()
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			ctx.Locals(field.Name, v.Field(i).Interface())
		}

		return ctx.Next()
	}
}

func extractToken(token string) (string, error) {
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return parts[1], nil
}

type ApiKeyMiddlewareConfig struct {
	Secret      string
	AuthKeyName string
}

func APIKeyMiddleware(config *ApiKeyMiddlewareConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get(config.AuthKeyName)
		if authHeader == "" {
			return exc.ForbiddenError("missing_api_key", "Missing API key")
		}
		if authHeader != config.Secret {
			return exc.ForbiddenError("invalid_api_key", "Wrong API key")
		}
		jwtLog.Info("Incoming request with %s: %s", config.AuthKeyName, authHeader)
		return c.Next()
	}
}
