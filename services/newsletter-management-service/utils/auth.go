package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func GetBearerToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	return strings.TrimPrefix(bearerToken, "Bearer ")
}

func ExtractSubFromToken(tokenString string) (string, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
		return "", fmt.Errorf("sub claim not found")
	}

	return "", fmt.Errorf("invalid token claims")
}
