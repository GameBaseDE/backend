package openapi

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strings"
)

func isValidLogin(request UserLogin, k kubernetesClient) (bool, error) {
	user, err := k.GetUserSecret(request.Email)
	if err != nil {
		return false, err
	}

	return user.Password == request.Password, nil
}

// Extract the authentication header from the request
// and check the authentication token against the valid authentication tokens
func isAuthorized(request *gin.Context) bool {
	token, err := ParseJwt(request)
	return err == nil && token.Valid
}

func ParseJwt(request *gin.Context) (*jwt.Token, error) {
	s := extractJwt(request)
	if s == "" {
		return nil, fmt.Errorf("invalid token")
	}

	token, err := jwt.ParseWithClaims(s, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return hmacSampleSecret(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
}

func extractJwt(request *gin.Context) string {
	authHeader := request.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	fields := strings.Fields(authHeader)
	if len(fields) == 2 && fields[0] == "Bearer" {
		return fields[1]
	}

	return ""
}

// Lookup the email address from the authentication token
func extractEmail(request *gin.Context) (string, error) {
	token, err := ParseJwt(request)
	if err == nil && token != nil {
		if claims, ok := token.Claims.(userClaims); ok && token.Valid {
			return claims.UserEmail, nil
		}
	}

	return "", err
}
