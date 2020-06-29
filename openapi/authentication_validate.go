package openapi

import (
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

// Lookup the email address from the authentication token
func extractEmail(request *gin.Context) string {
	authHeader := request.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	fields := strings.Fields(authHeader)
	if len(fields) == 2 && fields[0] == "Bearer" {
		removeExpiredTokens()
		for email, tokenDetails := range loggedInUsers {
			if tokenDetails.AccessToken == fields[1] {
				return email
			}
		}
	}

	return ""
}

func removeExpiredTokens() {
	for email, tokenDetails := range loggedInUsers {
		if tokenDetails.isExpired() && tokenDetails.isExpiredForever() {
			delete(loggedInUsers, email)
		}
	}
}

func removeToken(email string) {
	delete(loggedInUsers, email)
}
