package openapi

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func isValidUserLogin(request UserLogin) bool {
	dummy := UserLogin{
		Email:    "test@example.com",
		Password: "12345678",
	}

	return request.Email == dummy.Email && request.Password == dummy.Password
}

// Extract the authentication header from the request
// and check the authentication token against the valid authentication tokens
func isAuthorized(request *gin.Context) bool {
	authHeader := request.GetHeader("Authorization")
	if authHeader == "" {
		return false
	}

	fields := strings.Fields(authHeader)
	if len(fields) == 2 && fields[0] == "Bearer" {
		removeExpiredTokens()
		for email, tokenDetails := range loggedInUsers {
			if tokenDetails.AccessToken == fields[1] {
				return !tokenDetails.isExpired()
			}

			if tokenDetails.isExpired() && tokenDetails.isExpiredForever() {
				delete(loggedInUsers, email)
			}
		}
	}

	return false
}

func removeExpiredTokens() {
	for email, tokenDetails := range loggedInUsers {
		if tokenDetails.isExpired() && tokenDetails.isExpiredForever() {
			delete(loggedInUsers, email)
		}
	}
}
