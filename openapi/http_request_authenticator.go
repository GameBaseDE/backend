package openapi

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type httpRequestAuthenticator struct {
	nextHandler httpRequestHandler
}

func newHttpRequestAuthenticator() *httpRequestAuthenticator {
	return &httpRequestAuthenticator{nextHandler: newHttpRequestParser()}
}

// Login - Login a user and return a JWT with the user object
func (hr *httpRequestAuthenticator) Login(c *gin.Context) {
	var request UserLogin
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !isValidUserLogin(request) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := createToken(request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, User{
		Email:    "test@example.com",
		FullName: "Mr. Test",
		Token:    token.AccessToken,
	})
}

// Logout - Invalidate the passed JWT
func (hr *httpRequestAuthenticator) Logout(c *gin.Context) {
	if email := extractEmail(c); email != "" {
		removeToken(email)
		c.JSON(http.StatusOK, gin.H{"success": "success"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": "invalid authentication token"})
}

// Register - Register a user and return a JWT with the user object
func (hr *httpRequestAuthenticator) Register(c *gin.Context) {
	// TODO: authenticator.AuthRegisterPost(c)
}

// ListTemplates - Get a list of all available game server images
func (hr *httpRequestAuthenticator) ListTemplates(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.ListTemplates(c)
}

// GetStatus - Query status of all deployments
func (hr *httpRequestAuthenticator) GetStatus(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.GetStatus(c)
}

// ConfigureContainer - Configure a game server based on POST body
func (hr *httpRequestAuthenticator) ConfigureContainer(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.ConfigureContainer(c)
}

// DeployContainer - Deploy a game server based on POST body
func (hr *httpRequestAuthenticator) DeployContainer(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.DeployContainer(c)
}

// StartContainer - Start a game server/container
func (hr *httpRequestAuthenticator) StartContainer(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.StartContainer(c)
}

// StopContainer - Stop a game server/container
func (hr *httpRequestAuthenticator) StopContainer(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.StopContainer(c)
}

// RestartContainer - Restart a game server/container
func (hr *httpRequestAuthenticator) RestartContainer(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.RestartContainer(c)
}

// DeleteContainer - Delete deployment of game server
func (hr *httpRequestAuthenticator) DeleteContainer(c *gin.Context) {
	if !isAuthorized(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
		return
	}
	extractNamespace(c)
	hr.nextHandler.DeleteContainer(c)
}

// Sets the target namespace based on the Request JWT
func extractNamespace(c *gin.Context) {
	//FIXME could be combined with isAuthorized()
	c.Set("namespace", "gambaseprefix-testuser")
}
