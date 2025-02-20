package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/hewpao/hewpao-backend/config"
)

func JwtDecap(token string, cfg config.Config, ctx context.Context) (*string, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	claims := jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, err
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, err
	}

	return &email, nil
}
